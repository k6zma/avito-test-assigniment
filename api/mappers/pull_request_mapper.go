package mappers

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
)

func ToDomainPullRequest(dto dtos.PullRequest) (*models.PullRequest, error) {
	id, err := uuid.Parse(dto.PullRequestId)
	if err != nil {
		return nil, fmt.Errorf("invalid pull_request_id: %w", err)
	}

	authorID, err := uuid.Parse(dto.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("invalid author_id: %w", err)
	}

	status := valueobjects.PullRequestStatus(dto.Status)
	if status != valueobjects.OpenPullRequest && status != valueobjects.MergedPullRequest {
		return nil, fmt.Errorf("invalid status: %s", dto.Status)
	}

	reviewers := make([]uuid.UUID, 0, len(dto.AssignedReviewers))

	for _, rid := range dto.AssignedReviewers {
		uid, err := uuid.Parse(rid)
		if err != nil {
			return nil, fmt.Errorf("invalid reviewer uuid: %w", err)
		}

		reviewers = append(reviewers, uid)
	}

	pr := &models.PullRequest{
		ID:                id,
		Name:              dto.PullRequestName,
		AuthorID:          authorID,
		Status:            status,
		AssignedReviewers: reviewers,
	}

	if dto.CreatedAt != nil {
		pr.CreatedAt = dto.CreatedAt.UTC()
	} else {
		pr.CreatedAt = time.Time{}
	}

	if dto.MergedAt != nil {
		mergedAt := dto.MergedAt.UTC()
		pr.MergedAt = &mergedAt
	}

	if err = pr.Validate(); err != nil {
		return nil, err
	}

	return pr, nil
}

func ToDTOPullRequest(pr *models.PullRequest) dtos.PullRequest {
	dto := dtos.PullRequest{
		PullRequestId:     pr.ID.String(),
		PullRequestName:   pr.Name,
		AuthorId:          pr.AuthorID.String(),
		Status:            dtos.PullRequestStatus(pr.Status),
		AssignedReviewers: make([]string, 0, len(pr.AssignedReviewers)),
	}

	for _, r := range pr.AssignedReviewers {
		dto.AssignedReviewers = append(dto.AssignedReviewers, r.String())
	}

	createdAt := pr.CreatedAt
	dto.CreatedAt = &createdAt

	if pr.MergedAt != nil {
		mergedAt := pr.MergedAt.UTC()
		dto.MergedAt = &mergedAt
	}

	return dto
}
