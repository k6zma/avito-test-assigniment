package initialize

import (
	"context"
	"math/rand"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type BaseInitializeStrategy struct {
	numReviewers int
}

func NewBaseInitializeStrategy(n int) *BaseInitializeStrategy {
	if n <= 0 {
		n = 2
	}

	return &BaseInitializeStrategy{
		numReviewers: n,
	}
}

func (s *BaseInitializeStrategy) PickInitialReviewers(
	_ context.Context,
	_ *models.PullRequest,
	candidates []uuid.UUID,
) ([]uuid.UUID, error) {
	if len(candidates) == 0 {
		return []uuid.UUID{}, nil
	}

	if len(candidates) <= s.numReviewers {
		return candidates, nil
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[:s.numReviewers], nil
}
