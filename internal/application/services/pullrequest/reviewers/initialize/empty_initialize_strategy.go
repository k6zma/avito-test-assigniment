package initialize

import (
	"context"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type EmptyInitializeStrategy struct{}

func NewEmptyInitializeStrategy() *EmptyInitializeStrategy {
	return &EmptyInitializeStrategy{}
}

func (s *EmptyInitializeStrategy) PickInitialReviewers(
	_ context.Context,
	_ *models.PullRequest,
	_ []uuid.UUID,
) ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}
