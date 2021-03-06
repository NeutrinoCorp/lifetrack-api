package activity

import (
	"context"
	"strconv"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

// UpdateCommand requests a mutation of an aggregate.Activity
type UpdateCommand struct {
	Ctx context.Context
	ID  string

	CategoryID string
	Title      string
	Picture    string

	State string
}

// UpdateCommandHandler handles UpdateCommand(s) requests
type UpdateCommandHandler struct {
	repo repository.Activity
	bus  event.Bus
}

// NewUpdateCommandHandler creates a new UpdateCommandHandler
func NewUpdateCommandHandler(r repository.Activity, b event.Bus) *UpdateCommandHandler {
	return &UpdateCommandHandler{
		repo: r,
		bus:  b,
	}
}

func (h UpdateCommandHandler) Invoke(cmd UpdateCommand) error {
	isCmdEmpty := cmd.Title == "" && cmd.Picture == "" && cmd.State == "" && cmd.CategoryID == ""
	if cmd.ID == "" {
		return exception.NewRequiredField("activity_id")
	} else if isCmdEmpty {
		return exception.NewRequiredField("title, picture, state or category_id")
	}

	activity, _, err := h.repo.Fetch(cmd.Ctx, repository.ActivityCriteria{ID: cmd.ID})
	if err != nil {
		return err
	}
	snapshot := *activity[0]

	if err = h.updater(cmd, activity[0]); err != nil {
		return err
	}

	return h.persist(cmd.Ctx, activity[0], snapshot)
}

func (h UpdateCommandHandler) updater(cmd UpdateCommand, activity *aggregate.Activity) error {
	// update each field
	if cmd.CategoryID != "" {
		if err := activity.ChangeCategory(cmd.CategoryID); err != nil {
			return err
		}
	}
	if cmd.Title != "" {
		if err := activity.Rename(cmd.Title); err != nil {
			return err
		}
	}
	if cmd.Picture != "" {
		if err := activity.UploadPicture(cmd.Picture); err != nil {
			return err
		}
	}
	if cmd.State != "" {
		if err := h.updateState(cmd.State, activity); err != nil {
			return err
		}
	}

	return nil
}

func (h UpdateCommandHandler) updateState(stateStr string, activity *aggregate.Activity) error {
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		return exception.NewFieldFormat("state", "boolean")
	} else if state && !activity.State() {
		activity.Activate()
	} else if !state && activity.State() {
		activity.Deactivate()
	}

	return nil
}

func (h UpdateCommandHandler) persist(ctx context.Context, activity *aggregate.Activity,
	snapshot aggregate.Activity) error {
	if err := h.repo.Save(ctx, *activity); err != nil {
		return err
	}

	return h.pushEvents(ctx, activity, snapshot)
}

func (h UpdateCommandHandler) pushEvents(ctx context.Context, activity *aggregate.Activity,
	snapshot aggregate.Activity) error {
	if err := h.bus.Publish(ctx, activity.PullEvents()...); err != nil {
		// rollback
		if errR := h.repo.Save(ctx, snapshot); errR != nil {
			return errR
		}
		return err
	}

	return nil
}
