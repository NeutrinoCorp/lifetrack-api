package main

import (
	"context"
	"log"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"github.com/neutrinocorp/life-track-api/pkg/transport/categoryhandler"
)

var muxLambda *gorillamux.GorillaMuxAdapter
var cleaning func()

func init() {
	q, clean, err := dep.InjectGetCategoryQuery()
	if err != nil {
		log.Fatalf("failed to start get category query: %s", exception.GetDescription(err))
	}
	cleaning = clean

	h := categoryhandler.NewGet(q, mux.NewRouter().PathPrefix("/live").Subrouter())
	log.Print("category handler successfully started")
	muxLambda = gorillamux.New(h.GetRouter())
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.ProxyWithContext(ctx, req)
}

func main() {
	defer cleaning()
	lambda.Start(Handler)
}
