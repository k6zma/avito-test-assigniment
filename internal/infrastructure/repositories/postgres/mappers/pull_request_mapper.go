package mappers

import (
	"errors"
	"time"

	"github.com/google/uuid"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
)

func ToDomainPullRequest(
	pullRequestDB generated.PullRequest,
	reviewers []uuid.UUID,
) (*models.PullRequest, error) {
	var merged *time.Time

	if pullRequestDB.MergedAt.Valid {
		merged = &pullRequestDB.MergedAt.Time
	}

	pullRequestDomain := &models.PullRequest{
		ID:                UUIDFromPg(pullRequestDB.PullRequestID),
		Name:              pullRequestDB.PullRequestName,
		AuthorID:          UUIDFromPg(pullRequestDB.AuthorID),
		Status:            valueobjects.PullRequestStatus(pullRequestDB.Status),
		AssignedReviewers: reviewers,
		CreatedAt:         pullRequestDB.CreatedAt.Time,
		MergedAt:          merged,
	}

	if err := pullRequestDomain.Validate(); err != nil {
		return nil, errors.Join(apperrors.ErrValidation, err)
	}

	return pullRequestDomain, nil
}

func ToDBCreatePullRequest(
	pullRequestDomain *models.PullRequest,
) generated.CreatePullRequestParams {
	return generated.CreatePullRequestParams{
		PullRequestID:   UUIDToPg(pullRequestDomain.ID),
		PullRequestName: pullRequestDomain.Name,
		AuthorID:        UUIDToPg(pullRequestDomain.AuthorID),
	}
}
