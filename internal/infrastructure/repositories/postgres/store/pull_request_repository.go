package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/mappers"
	"github.com/k6zma/avito-test-assigniment/pkg/logger"
)

type PullRequestRepository struct {
	queries *generated.Queries
	logger  *slog.Logger
}

func NewPullRequestRepository(
	q *generated.Queries,
	logger *slog.Logger,
) repositories.PullRequestRepository {
	const pullRequestEntity = "pull_request"

	return &PullRequestRepository{
		queries: q,
		logger:  logger.WithGroup(pullRequestEntity).WithGroup(postgresDBType),
	}
}

func (r *PullRequestRepository) Create(
	ctx context.Context,
	pullRequest *models.PullRequest,
) (err error) {
	const createMethod = "Create"

	log := r.logger.WithGroup(createMethod)

	operationLog := logger.StartOperation(
		log,
		"creating pull request in store",
		slog.String("pull_request_id", pullRequest.ID.String()),
		slog.String("pull_request_name", pullRequest.Name),
		slog.String("author_id", pullRequest.AuthorID.String()),
		slog.Any("status", pullRequest.Status),
	)

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.String("pull_request_id", pullRequest.ID.String()),
			slog.Int("reviewers_count", len(pullRequest.AssignedReviewers)),
		)
	}()

	params := mappers.ToDBCreatePullRequest(pullRequest)

	if err := r.queries.CreatePullRequest(ctx, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.Join(apperrors.ErrPRExists, err)
		}

		return fmt.Errorf("failed create pull request in postgres storage: %w", err)
	}

	for _, reviewerID := range pullRequest.AssignedReviewers {
		if err := r.queries.AddReviewer(ctx, mappers.ToDBReviewer(pullRequest.ID, reviewerID)); err != nil {
			return fmt.Errorf("failed add reviewer to pull request in postgres storage: %w", err)
		}
	}

	return nil
}

func (r *PullRequestRepository) GetByID(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (pullRequest *models.PullRequest, err error) {
	const getByIDMethod = "GetByID"

	log := r.logger.WithGroup(getByIDMethod)

	operationLog := logger.StartOperation(
		log,
		"getting pull request from store",
		slog.String("pull_request_id", pullRequestID.String()),
	)

	defer func() {
		var extra []any

		if pullRequest != nil {
			extra = append(
				extra,
				slog.String("pull_request_id", pullRequest.ID.String()),
				slog.Any("status", pullRequest.Status),
				slog.Int("reviewers_count", len(pullRequest.AssignedReviewers)),
			)
		}

		operationLog.FinishOperation(
			&err,
			append([]any{slog.String("pull_request_id", pullRequestID.String())}, extra...)...)
	}()

	dbPullRequest, err := r.queries.GetPullRequest(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrPRNotFound, err)
		}

		return nil, fmt.Errorf("failed get pull request from postgres storage: %w", err)
	}

	reviewerIDS, err := r.queries.ListReviewers(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		return nil, fmt.Errorf("failed get pull request reviewers from postgres storage: %w", err)
	}

	reviewers := make([]uuid.UUID, 0, len(reviewerIDS))
	for _, reviewerID := range reviewerIDS {
		reviewers = append(reviewers, mappers.ToDomainReviewer(reviewerID))
	}

	pullRequest, err = mappers.ToDomainPullRequest(dbPullRequest, reviewers)
	if err != nil {
		return nil, fmt.Errorf("failed convert pull request to domain model: %w", err)
	}

	return pullRequest, nil
}

func (r *PullRequestRepository) Merge(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (pullRequest *models.PullRequest, err error) {
	const mergeMethod = "Merge"

	log := r.logger.WithGroup(mergeMethod)

	operationLog := logger.StartOperation(
		log,
		"merging pull request in store",
		slog.String("pull_request_id", pullRequestID.String()),
	)

	defer func() {
		var extra []any

		if pullRequest != nil {
			extra = append(
				extra,
				slog.Any("status", pullRequest.Status),
				slog.Int("reviewers_count", len(pullRequest.AssignedReviewers)),
			)
		}

		operationLog.FinishOperation(
			&err,
			append([]any{slog.String("pull_request_id", pullRequestID.String())}, extra...)...,
		)
	}()

	dbPullRequest, err := r.queries.MergePullRequest(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrPRNotFound, err)
		}

		return nil, fmt.Errorf("failed merge pull request in postgres storage: %w", err)
	}

	reviewerIDS, err := r.queries.ListReviewers(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		return nil, fmt.Errorf("failed get pull request reviewers from postgres storage: %w", err)
	}

	reviewers := make([]uuid.UUID, 0, len(reviewerIDS))

	for _, reviewerID := range reviewerIDS {
		reviewers = append(reviewers, mappers.ToDomainReviewer(reviewerID))
	}

	pullRequest, err = mappers.ToDomainPullRequest(dbPullRequest, reviewers)
	if err != nil {
		return nil, fmt.Errorf("failed convert pull request to domain model: %w", err)
	}

	return pullRequest, nil
}

func (r *PullRequestRepository) ListAssignedToUser(
	ctx context.Context,
	reviewerID uuid.UUID,
) (pullRequests []*models.PullRequest, err error) {
	const listAssignedMethod = "ListAssignedToUser"

	log := r.logger.WithGroup(listAssignedMethod)

	operationLog := logger.StartOperation(
		log,
		"listing assigned pull requests in store",
		slog.String("reviewer_id", reviewerID.String()),
	)

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.String("reviewer_id", reviewerID.String()),
			slog.Int("pull_requests_count", len(pullRequests)),
		)
	}()

	dbPullRequests, err := r.queries.ListPRsAssignedToUser(ctx, mappers.UUIDToPg(reviewerID))
	if err != nil {
		return nil, fmt.Errorf("failed get list pull requests from postgres storage: %w", err)
	}

	pullRequests = make([]*models.PullRequest, 0, len(dbPullRequests))

	for _, dbPullRequest := range dbPullRequests {
		reviewerIDS, err := r.queries.ListReviewers(ctx, dbPullRequest.PullRequestID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed get pull request reviewers from postgres storage: %w",
				err,
			)
		}

		reviewers := make([]uuid.UUID, 0, len(reviewerIDS))
		for _, reviewerID := range reviewerIDS {
			reviewers = append(reviewers, mappers.ToDomainReviewer(reviewerID))
		}

		pullRequest, err := mappers.ToDomainPullRequest(dbPullRequest, reviewers)
		if err != nil {
			return nil, fmt.Errorf("failed convert pull request to domain model: %w", err)
		}

		pullRequests = append(pullRequests, pullRequest)
	}

	return pullRequests, nil
}

func (r *PullRequestRepository) AddReviewer(
	ctx context.Context,
	pullRequestID, reviewerID uuid.UUID,
) (err error) {
	const addReviewerMethod = "AddReviewer"

	log := r.logger.WithGroup(addReviewerMethod)

	operationLog := logger.StartOperation(
		log,
		"adding reviewer in store",
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)

	params := mappers.ToDBReviewer(pullRequestID, reviewerID)

	if err = r.queries.AddReviewer(ctx, params); err != nil {
		return fmt.Errorf("failed add reviewer to pull request in postgres storage: %w", err)
	}

	return nil
}

func (r *PullRequestRepository) RemoveReviewer(
	ctx context.Context,
	pullRequestID, reviewerID uuid.UUID,
) (err error) {
	const removeReviewerMethod = "RemoveReviewer"

	log := r.logger.WithGroup(removeReviewerMethod)

	operationLog := logger.StartOperation(
		log,
		"removing reviewer in store",
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)

	if err = r.queries.RemoveReviewer(ctx, generated.RemoveReviewerParams{
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
) (reviewerIDS []uuid.UUID, err error) {
	const listReviewersMethod = "ListReviewers"

	log := r.logger.WithGroup(listReviewersMethod)

	operationLog := logger.StartOperation(
		log,
		"listing reviewers in store",
		slog.String("pull_request_id", pullRequestID.String()),
	)

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.String("pull_request_id", pullRequestID.String()),
			slog.Int("reviewers_count", len(reviewerIDS)),
		)
	}()

	reviewerIDSFromDB, err := r.queries.ListReviewers(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		return nil, fmt.Errorf("failed get pull request reviewers from postgres storage: %w", err)
	}

	reviewerIDS = make([]uuid.UUID, 0, len(reviewerIDSFromDB))
	for _, reviewerID := range reviewerIDSFromDB {
		reviewerIDS = append(reviewerIDS, mappers.ToDomainReviewer(reviewerID))
	}

	return reviewerIDS, nil
}

func (r *PullRequestRepository) CountReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (count int64, err error) {
	const countReviewersMethod = "CountReviewers"

	log := r.logger.WithGroup(countReviewersMethod)

	operationLog := logger.StartOperation(
		log,
		"counting reviewers in store",
		slog.String("pull_request_id", pullRequestID.String()),
	)
	defer operationLog.FinishOperation(&err, slog.String("pull_request_id", pullRequestID.String()))

	count, err = r.queries.CountReviewers(ctx, mappers.UUIDToPg(pullRequestID))
	if err != nil {
		return 0, fmt.Errorf("failed count pull request reviewers in postgres storage: %w", err)
	}

	return count, nil
}

func (r *PullRequestRepository) ActiveCandidatesForReassign(
	ctx context.Context,
	userID uuid.UUID,
) (userIDs []uuid.UUID, err error) {
	const activeCandidatesMethod = "ActiveCandidatesForReassign"

	log := r.logger.WithGroup(activeCandidatesMethod)

	operationLog := logger.StartOperation(
		log,
		"listing active reassign candidates in store",
		slog.String("user_id", userID.String()),
	)

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.String("user_id", userID.String()),
			slog.Int("candidates_count", len(userIDs)),
		)
	}()

	userIDSFromDB, err := r.queries.ActiveCandidatesForReassign(ctx, mappers.UUIDToPg(userID))
	if err != nil {
		return nil, fmt.Errorf(
			"failed get active candidates for reassign from postgres storage: %w",
			err,
		)
	}

	userIDs = make([]uuid.UUID, 0, len(userIDSFromDB))

	for _, userIDFromDB := range userIDSFromDB {
		userIDs = append(userIDs, mappers.UUIDFromPg(userIDFromDB))
	}

	return userIDs, nil
}
