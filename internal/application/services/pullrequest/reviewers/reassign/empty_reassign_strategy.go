package reassign

import (
	"context"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type EmptyReassignStrategy struct{}

func NewEmptyReassignStrategy() *EmptyReassignStrategy {
	return &EmptyReassignStrategy{}
}

func (s *EmptyReassignStrategy) PickReplacementReviewers(
	_ context.Context,
	_ *models.PullRequest,
	_ uuid.UUID,
	_ []uuid.UUID,
) (uuid.UUID, error) {
	return uuid.Nil, nil
}
