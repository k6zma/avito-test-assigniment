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
	pullRequestID, err := uuid.Parse(dto.PullRequestId)
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

	for _, reviewerID := range dto.AssignedReviewers {
		userID, err := uuid.Parse(reviewerID)
		if err != nil {
			return nil, fmt.Errorf("invalid reviewer uuid: %w", err)
		}

		reviewers = append(reviewers, userID)
	}

	pullRequest := &models.PullRequest{
		ID:                pullRequestID,
		Name:              dto.PullRequestName,
		AuthorID:          authorID,
		Status:            status,
		AssignedReviewers: reviewers,
	}

	if dto.CreatedAt != nil {
		pullRequest.CreatedAt = dto.CreatedAt.UTC()
	} else {
		pullRequest.CreatedAt = time.Time{}
	}

	if dto.MergedAt != nil {
		mergedAt := dto.MergedAt.UTC()
		pullRequest.MergedAt = &mergedAt
	}

	if err = pullRequest.Validate(); err != nil {
		return nil, err
	}

	return pullRequest, nil
}

func ToDTOPullRequest(pullRequest *models.PullRequest) dtos.PullRequest {
	pullRequestDTO := dtos.PullRequest{
		PullRequestId:     pullRequest.ID.String(),
		PullRequestName:   pullRequest.Name,
		AuthorId:          pullRequest.AuthorID.String(),
		Status:            dtos.PullRequestStatus(pullRequest.Status),
		AssignedReviewers: make([]string, 0, len(pullRequest.AssignedReviewers)),
	}

	for _, reviewerID := range pullRequest.AssignedReviewers {
		pullRequestDTO.AssignedReviewers = append(
			pullRequestDTO.AssignedReviewers,
			reviewerID.String(),
		)
	}

	createdAt := pullRequest.CreatedAt
	pullRequestDTO.CreatedAt = &createdAt

	if pullRequest.MergedAt != nil {
		mergedAt := pullRequest.MergedAt.UTC()
		pullRequestDTO.MergedAt = &mergedAt
	}

	return pullRequestDTO
}
