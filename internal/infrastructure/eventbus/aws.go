package eventbus

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"strings"
	"sync"
)

// AWSEventBus is the event.Bus implementation using AWS SNS and SQS
type AWSEventBus struct {
	sess *session.Session
	cfg  infrastructure.Configuration
	mu   *sync.Mutex
}

// NewAWSEventBus creates a concrete struct of AWSEventBus
func NewAWSEventBus(s *session.Session, cfg infrastructure.Configuration) *AWSEventBus {
	return &AWSEventBus{
		sess: s,
		cfg:  cfg,
		mu:   new(sync.Mutex),
	}
}

func (b *AWSEventBus) Publish(ctx context.Context, e ...event.Domain) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(e) == 0 {
		return exception.NewRequiredField("domain event")
	}

	svc := NewSNSConn(b.sess, b.cfg.Category.Event.Region)
	for _, ev := range e {
		ev.TopicToUnderscore()
		// Get topic Arn before publish
		topicArn, err := b.getTopicArn(ctx, svc, ev.Topic)
		if err != nil {
			return err
		}

		eventJSON, err := ev.MarshalBinary()
		if err != nil {
			return err
		}

		_, err = svc.PublishWithContext(ctx, &sns.PublishInput{
			Message:          aws.String(string(eventJSON)),
			MessageStructure: aws.String("json"),
			TopicArn:         aws.String(topicArn),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case sns.ErrCodeResourceNotFoundException:
					return exception.NewNotFound(ev.Topic)
				case sns.ErrCodeInvalidParameterException:
					return exception.NewFieldFormat(ev.Topic+" parameter", "valid topic parameter ")
				case sns.ErrCodeInvalidParameterValueException:
					return exception.NewFieldFormat(ev.Topic+" parameter", "valid topic parameter value")
				case sns.ErrCodeThrottledException:
					return exception.NewNetworkCall("aws sns topic "+ev.Topic, b.cfg.Category.Event.Region)
				}
			}

			return err
		}
	}

	return nil
}

func (b *AWSEventBus) SubscribeTo(ctx context.Context, t event.Topic) (*event.Domain, error) {
	svc := NewSQSConn(b.sess, b.cfg.Category.Event.Region)

	queueURL, err := b.getQueueURL(ctx, svc, string(t))
	if err != nil {
		return nil, err
	}

	// Long-polling strategy
	o, err := svc.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: aws.StringSlice([]string{
			"SentTimestamp",
		}),
		MaxNumberOfMessages: aws.Int64(1),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
		WaitTimeSeconds: aws.Int64(20),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sqs.ErrCodeOverLimit:
				return nil, exception.NewNetworkCall("aws sns topic "+string(t), b.cfg.Category.Event.Region)
			}
		}

		return nil, err
	}

	// Use adapter func
	e, err := b.getDomainEvent(o.Messages[0])
	if err != nil {
		return nil, err
	}

	return e, nil
}

// getTopicArn returns the given topic ARN from given session's AWS resources
func (b AWSEventBus) getTopicArn(ctx context.Context, svc *sns.SNS, topic string) (string, error) {
	nextToken := ""

	for {
		result, err := svc.ListTopicsWithContext(ctx, &sns.ListTopicsInput{
			NextToken: aws.String(nextToken),
		})
		if err != nil {
			return "", exception.NewNotFound("topics")
		}

		// Up to 100 topics
		for _, t := range result.Topics {
			// Search for given topic
			spl := strings.Split(*t.TopicArn, ":")
			if spl[len(spl)-1] == topic {
				return *t.TopicArn, nil
			}
		}

		// If no more to fetch, then break
		if result.NextToken == nil || *result.NextToken == "" {
			break
		}

		nextToken = *result.NextToken
	}

	return "", exception.NewNotFound("topic")
}

// getQueueURL returns a queue URL for the given queue name
func (b AWSEventBus) getQueueURL(ctx context.Context, svc *sqs.SQS, queue string) (string, error) {
	result, err := svc.GetQueueUrlWithContext(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	})
	if err != nil || result.QueueUrl == nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sqs.ErrCodeQueueDoesNotExist:
				return "", exception.NewNotFound("queue " + queue)
			case sqs.ErrCodeOverLimit:
				return "", exception.NewNetworkCall("aws sqs queue "+queue, b.cfg.Category.Event.Region)
			}
		}

		return "", exception.NewNotFound("topic queue")
	}

	return *result.QueueUrl, nil
}

// getDomainEvent adapts sqs.Message into a domain event
func (b AWSEventBus) getDomainEvent(msg *sqs.Message) (*event.Domain, error) {
	if msg == nil || msg.Body == nil || msg.ReceiptHandle == nil {
		return nil, exception.NewRequiredField("message")
	}

	e := new(event.Domain)
	if err := e.UnmarshalBinary([]byte(*msg.Body)); err != nil {
		return nil, err
	}
	e.Acknowledge = *msg.ReceiptHandle

	return e, nil
}
