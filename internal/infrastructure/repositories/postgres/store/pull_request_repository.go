package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/mappers"
)

type PullRequestRepository struct {
	queries *generated.Queries
}

func NewPullRequestRepository(q *generated.Queries) repositories.PullRequestRepository {
	return &PullRequestRepository{queries: q}
}

func (r *PullRequestRepository) Create(ctx context.Context, pr *models.PullRequest) error {
	params := mappers.ToDBCreatePullRequest(pr)

	if err := r.queries.CreatePullRequest(ctx, params); err != nil {
		return fmt.Errorf("failed create pull request in postgres storage: %w", err)
	}

	for _, reviewerID := range pr.AssignedReviewers {
		if err := r.queries.AddReviewer(ctx, mappers.ToDBReviewer(pr.ID, reviewerID)); err != nil {
			return fmt.Errorf("failed add reviewer to pull request in postgres storage: %w", err)
		}
	}

	return nil
}

func (r *PullRequestRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.PullRequest, error) {
	dbPR, err := r.queries.GetPullRequest(ctx, mappers.UUIDToPg(id))
	if err != nil {
		return nil, fmt.Errorf("failed get pull request from postgres storage: %w", err)
	}

	reviewerIDs, err := r.queries.ListReviewers(ctx, mappers.UUIDToPg(id))
	if err != nil {
		return nil, fmt.Errorf("failed get pull request reviewers from postgres storage: %w", err)
	}

	reviewers := make([]uuid.UUID, 0, len(reviewerIDs))
	for _, rid := range reviewerIDs {
		reviewers = append(reviewers, mappers.ToDomainReviewer(rid))
	}

	domainPR, err := mappers.ToDomainPullRequest(dbPR, reviewers)
	if err != nil {
		return nil, fmt.Errorf("failed convert pull request to domain model: %w", err)
	}

	return domainPR, nil
}

func (r *PullRequestRepository) Merge(
	ctx context.Context,
	id uuid.UUID,
) (*models.PullRequest, error) {
	dbPR, err := r.queries.MergePullRequest(ctx, mappers.UUIDToPg(id))
	if err != nil {
		return nil, fmt.Errorf("failed merge pull request in postgres storage: %w", err)
	}

	reviewerIDs, err := r.queries.ListReviewers(ctx, mappers.UUIDToPg(id))
	if err != nil {
		return nil, fmt.Errorf("failed get pull request reviewers from postgres storage: %w", err)
	}

	reviewers := make([]uuid.UUID, 0, len(reviewerIDs))
	for _, rid := range reviewerIDs {
		reviewers = append(reviewers, mappers.ToDomainReviewer(rid))
	}

	domainPR, err := mappers.ToDomainPullRequest(dbPR, reviewers)
	if err != nil {
		return nil, fmt.Errorf("failed convert pull request to domain model: %w", err)
	}

	return domainPR, nil
}

func (r *PullRequestRepository) ListAssignedToUser(
	ctx context.Context,
	reviewerID uuid.UUID,
) ([]*models.PullRequest, error) {
	dbPRs, err := r.queries.ListPRsAssignedToUser(ctx, mappers.UUIDToPg(reviewerID))
	if err != nil {
		return nil, fmt.Errorf("failed get list pull requests from postgres storage: %w", err)
	}

	result := make([]*models.PullRequest, 0, len(dbPRs))

	for _, dbPR := range dbPRs {
		reviewerIDs, err := r.queries.ListReviewers(ctx, dbPR.PullRequestID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed get pull request reviewers from postgres storage: %w",
				err,
			)
		}

		reviewers := make([]uuid.UUID, 0, len(reviewerIDs))
		for _, rid := range reviewerIDs {
			reviewers = append(reviewers, mappers.ToDomainReviewer(rid))
		}

		pr, err := mappers.ToDomainPullRequest(dbPR, reviewers)
		if err != nil {
			return nil, fmt.Errorf("failed convert pull request to domain model: %w", err)
		}

		result = append(result, pr)
	}

	return result, nil
}

func (r *PullRequestRepository) AddReviewer(
	ctx context.Context,
	pullRequestID, reviewerID uuid.UUID,
) error {
	params := mappers.ToDBReviewer(pullRequestID, reviewerID)

	if err := r.queries.AddReviewer(ctx, params); err != nil {
		return fmt.Errorf("failed add reviewer to pull request in postgres storage: %w", err)
	}

	return nil
}

func (r *PullRequestRepository) RemoveReviewer(
	ctx context.Context,
	pullRequestID, reviewerID uuid.UUID,
) error {
	if err := r.queries.RemoveReviewer(ctx, generated.RemoveReviewerParams{
		PullRequestID: mappers.UUIDToPg(pullRequestID),
		ReviewerID:    mappers.UUIDToPg(reviewerID),
	}); err != nil {
		return fmt.Errorf("failed remove reviewer from pull request in postgres storage: %w", err)
	}

	return nil
}

func (r *PullRequestRepository) ListReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
) ([]uuid.UUID, error) {
	ids, err := r.queries.ListReviewers(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		return nil, fmt.Errorf("failed get pull request reviewers from postgres storage: %w", err)
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		result = append(result, mappers.ToDomainReviewer(id))
	}

	return result, nil
}

func (r *PullRequestRepository) CountReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (int64, error) {
	count, err := r.queries.CountReviewers(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		return 0, fmt.Errorf("failed count pull request reviewers in postgres storage: %w", err)
	}

	return count, nil
}

func (r *PullRequestRepository) ActiveCandidatesForReassign(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {
	ids, err := r.queries.ActiveCandidatesForReassign(ctx, mappers.UUIDToPg(userID))
	if err != nil {
		return nil, fmt.Errorf(
			"failed get active candidates for reassign from postgres storage: %w",
			err,
		)
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		result = append(result, mappers.UUIDFromPg(id))
	}

	return result, nil
}
