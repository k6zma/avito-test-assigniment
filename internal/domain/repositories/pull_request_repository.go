package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PullRequest, error)
	Merge(ctx context.Context, id uuid.UUID) (*models.PullRequest, error)
	ListAssignedToUser(ctx context.Context, reviewerID uuid.UUID) ([]*models.PullRequest, error)
	AddReviewer(ctx context.Context, pullRequestID, reviewerID uuid.UUID) error
	RemoveReviewer(ctx context.Context, pullRequestID, reviewerID uuid.UUID) error
	ListReviewers(ctx context.Context, pullRequestID uuid.UUID) ([]uuid.UUID, error)
	CountReviewers(ctx context.Context, pullRequestID uuid.UUID) (int64, error)
	ActiveCandidatesForReassign(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}
