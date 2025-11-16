package reassign

import (
	"context"
	"math/rand"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type BaseReassignStrategy struct{}

func NewBaseReassignStrategy() *BaseReassignStrategy {
	return &BaseReassignStrategy{}
}

func (s *BaseReassignStrategy) PickReplacementReviewers(
	_ context.Context,
	pr *models.PullRequest,
	leavingReviewer uuid.UUID,
	candidates []uuid.UUID,
) (uuid.UUID, error) {
	filtered := make([]uuid.UUID, 0, len(candidates))

	for _, id := range candidates {
		if id == leavingReviewer || id == pr.AuthorID {
			continue
		}

		filtered = append(filtered, id)
	}

	if len(filtered) == 0 {
		return uuid.Nil, nil
	}

	idx := rand.Intn(len(filtered))

	return filtered[idx], nil
}
