package initialize

import (
	"context"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type ReviewerInitializeStrategy interface {
	PickInitialReviewers(
		ctx context.Context,
		pullRequest *models.PullRequest,
		candidates []uuid.UUID,
	) ([]uuid.UUID, error)
}
