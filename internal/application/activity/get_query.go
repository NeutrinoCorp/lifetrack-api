package activity

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

// GetQuery requests a single model.Activity
type GetQuery struct {
	repo repository.Activity
}

// NewGetQuery creates a GetQuery
func NewGetQuery(r repository.Activity) *GetQuery {
	return &GetQuery{repo: r}
}

func (q GetQuery) Query(ctx context.Context, id string) (*model.Activity, error) {
	acts, _, err := q.repo.Fetch(ctx, repository.ActivityCriteria{ID: id})
	if err != nil {
		return nil, err
	}

	return acts[0].MarshalPrimitive(), nil
}
