package mappers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
)

func ToDomainReviewer(id pgtype.UUID) uuid.UUID {
	return UUIDFromPg(id)
}

func ToDBReviewer(pullRequestID uuid.UUID, reviewerID uuid.UUID) generated.AddReviewerParams {
	return generated.AddReviewerParams{
		ID:            UUIDToPg(uuid.New()),
		PullRequestID: UUIDToPg(pullRequestID),
		ReviewerID:    UUIDToPg(reviewerID),
	}
}
