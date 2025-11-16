package reassign

import (
	"context"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type ReviewerReassignStrategy interface {
	PickReplacementReviewers(
		ctx context.Context,
		pr *models.PullRequest,
		leavingReviewer uuid.UUID,
		candidates []uuid.UUID,
	) (uuid.UUID, error)
}
