package pullrequest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest/reviewers/initialize"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest/reviewers/reassign"
	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
	"github.com/k6zma/avito-test-assigniment/pkg/logger"
)

type PullRequestService struct {
	pullRequests     repositories.PullRequestRepository
	users            repositories.UserRepository
	initStrategy     initialize.ReviewerInitializeStrategy
	reassignStrategy reassign.ReviewerReassignStrategy
	logger           *slog.Logger
}

func NewPullRequestService(
	prRepo repositories.PullRequestRepository,
	userRepo repositories.UserRepository,
	initStrategy initialize.ReviewerInitializeStrategy,
	reassignStrategy reassign.ReviewerReassignStrategy,
	logger *slog.Logger,
) *PullRequestService {
	const pullRequestEntity = "pull_request"

	if initStrategy == nil {
		initStrategy = initialize.NewBaseInitializeStrategy(2)
	}

	if reassignStrategy == nil {
		reassignStrategy = reassign.NewBaseReassignStrategy()
	}

	return &PullRequestService{
		pullRequests:     prRepo,
		users:            userRepo,
		initStrategy:     initStrategy,
		reassignStrategy: reassignStrategy,
		logger:           logger.WithGroup(pullRequestEntity),
	}
}

func (s *PullRequestService) CreatePullRequest(
	ctx context.Context,
	pullRequestID uuid.UUID,
	name string,
	authorID uuid.UUID,
) (pullRequest *models.PullRequest, err error) {
	const createPullRequestMethod = "CreatePullRequest"

	log := s.logger.WithGroup(createPullRequestMethod)

	operationLog := logger.StartOperation(
		log,
		"creating pull request",
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("author_id", authorID.String()),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("author_id", authorID.String()),
	)

	author, err := s.users.GetByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed get author user: %w", err)
	}

	if !author.IsActive {
		return nil, errors.Join(
			apperrors.ErrUserInactive,
			fmt.Errorf("author %s is not active", authorID),
		)
	}

	teamActive, err := s.users.ListActiveUsersInTeam(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed list active users in team: %w", err)
	}

	candidates := make([]uuid.UUID, 0, len(teamActive))

	for _, userID := range teamActive {
		if userID != authorID {
			candidates = append(candidates, userID)
		}
	}

	pullRequest, err = models.NewPullRequest(
		pullRequestID,
		name,
		authorID,
		valueobjects.OpenPullRequest,
		[]uuid.UUID{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed create pull request domain model: %w", err)
	}

	assigned, err := s.initStrategy.PickInitialReviewers(ctx, pullRequest, candidates)
	if err != nil {
		return nil, fmt.Errorf("failed pick initial reviewers: %w", err)
	}

	pullRequest.AssignedReviewers = assigned

	if err = s.pullRequests.Create(ctx, pullRequest); err != nil {
		return nil, fmt.Errorf("failed create pull request: %w", err)
	}

	return pullRequest, nil
}

func (s *PullRequestService) GetPullRequest(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (pullRequest *models.PullRequest, err error) {
	const getPullRequestMethod = "GetPullRequest"

	log := s.logger.WithGroup(getPullRequestMethod)

	operationLog := logger.StartOperation(
		log,
		"getting pull request",
		slog.String("pull_request_id", pullRequestID.String()),
	)
	defer operationLog.FinishOperation(&err, slog.String("pull_request_id", pullRequestID.String()))

	pullRequest, err = s.pullRequests.GetByID(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed get pull request by id: %w", err)
	}

	return pullRequest, nil
}

func (s *PullRequestService) MergePullRequest(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (pullRequest *models.PullRequest, err error) {
	const mergePullRequestMethod = "MergePullRequest"

	log := s.logger.WithGroup(mergePullRequestMethod)

	operationLog := logger.StartOperation(
		log,
		"merging pull request",
		slog.String("pull_request_id", pullRequestID.String()),
	)

	defer operationLog.FinishOperation(&err, slog.String("pull_request_id", pullRequestID.String()))

	pullRequest, err = s.pullRequests.Merge(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed merge pull request: %w", err)
	}

	return pullRequest, nil
}

func (s *PullRequestService) ListPullRequestsAssignedToUser(
	ctx context.Context,
	reviewerID uuid.UUID,
) (pullRequests []*models.PullRequest, err error) {
	const listAssignedMethod = "ListPullRequestsAssignedToUser"

	log := s.logger.WithGroup(listAssignedMethod)

	operationLog := logger.StartOperation(
		log,
		"listing pull requests assigned to user",
		slog.String("reviewer_id", reviewerID.String()),
	)
	defer operationLog.FinishOperation(&err, slog.String("reviewer_id", reviewerID.String()))

	pullRequests, err = s.pullRequests.ListAssignedToUser(ctx, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed list pull requests assigned to user: %w", err)
	}

	return pullRequests, nil
}

func (s *PullRequestService) AddReviewer(
	ctx context.Context,
	pullRequestID, reviewerID uuid.UUID,
) (err error) {
	const addReviewerMethod = "AddReviewer"

	log := s.logger.WithGroup(addReviewerMethod)

	operationLog := logger.StartOperation(
		log,
		"adding reviewer",
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)

	pullRequest, err := s.pullRequests.GetByID(ctx, pullRequestID)
	if err != nil {
		return fmt.Errorf("failed get pull request: %w", err)
	}

	if pullRequest.Status != valueobjects.OpenPullRequest {
		return fmt.Errorf("cannot modify reviewers for non-open pull request")
	}

	user, err := s.users.GetByID(ctx, reviewerID)
	if err != nil {
		return fmt.Errorf("failed get reviewer user: %w", err)
	}

	if !user.IsActive {
		return fmt.Errorf("reviewer %s is not active", reviewerID)
	}

	existing, err := s.pullRequests.ListReviewers(ctx, pullRequestID)
	if err != nil {
		return fmt.Errorf("failed list existing reviewers: %w", err)
	}

	for _, id := range existing {
		if id == reviewerID {
			return nil
		}
	}

	if err = s.pullRequests.AddReviewer(ctx, pullRequestID, reviewerID); err != nil {
		return fmt.Errorf("failed add reviewer to pull request: %w", err)
	}

	return nil
}

func (s *PullRequestService) RemoveReviewer(
	ctx context.Context,
	pullRequestID, reviewerID uuid.UUID,
) (err error) {
	const removeReviewerMethod = "RemoveReviewer"

	log := s.logger.WithGroup(removeReviewerMethod)

	operationLog := logger.StartOperation(
		log,
		"removing reviewer",
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("reviewer_id", reviewerID.String()),
	)

	pullRequest, err := s.pullRequests.GetByID(ctx, pullRequestID)
	if err != nil {
		return fmt.Errorf("failed get pull request: %w", err)
	}

	if pullRequest.Status != valueobjects.OpenPullRequest {
		return fmt.Errorf("cannot modify reviewers for non-open pull request")
	}

	if err = s.pullRequests.RemoveReviewer(ctx, pullRequestID, reviewerID); err != nil {
		return fmt.Errorf("failed remove reviewer from pull request: %w", err)
	}

	return nil
}

func (s *PullRequestService) ListReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (reviewers []uuid.UUID, err error) {
	const listReviewersMethod = "ListReviewers"

	log := s.logger.WithGroup(listReviewersMethod)

	operationLog := logger.StartOperation(
		log,
		"listing reviewers",
		slog.String("pull_request_id", pullRequestID.String()),
	)
	defer operationLog.FinishOperation(&err, slog.String("pull_request_id", pullRequestID.String()))

	reviewers, err = s.pullRequests.ListReviewers(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed list reviewers for pull request: %w", err)
	}

	return reviewers, nil
}

func (s *PullRequestService) CountReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (count int64, err error) {
	const countReviewersMethod = "CountReviewers"

	log := s.logger.WithGroup(countReviewersMethod)

	operationLog := logger.StartOperation(
		log,
		"counting reviewers",
		slog.String("pull_request_id", pullRequestID.String()),
	)
	defer operationLog.FinishOperation(&err, slog.String("pull_request_id", pullRequestID.String()))

	count, err = s.pullRequests.CountReviewers(ctx, pullRequestID)
	if err != nil {
		return 0, fmt.Errorf("failed count reviewers for pull request: %w", err)
	}

	return count, nil
}

func (s *PullRequestService) ReassignReviewer(
	ctx context.Context,
	pullRequestID, leavingReviewerID uuid.UUID,
) (err error) {
	const reassignReviewerMethod = "ReassignReviewer"

	log := s.logger.WithGroup(reassignReviewerMethod)

	operationLog := logger.StartOperation(
		log,
		"reassigning reviewer",
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("leaving_reviewer_id", leavingReviewerID.String()),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("pull_request_id", pullRequestID.String()),
		slog.String("leaving_reviewer_id", leavingReviewerID.String()),
	)

	pullRequest, err := s.pullRequests.GetByID(ctx, pullRequestID)
	if err != nil {
		return err
	}

	if pullRequest.Status != valueobjects.OpenPullRequest {
		return apperrors.ErrPRNotOpen
	}

	found := slices.Contains(pullRequest.AssignedReviewers, leavingReviewerID)

	if !found {
		return errors.Join(
			apperrors.ErrNotAssigned,
			fmt.Errorf("user %s is not a reviewer of PR %s", leavingReviewerID, pullRequestID),
		)
	}

	leavingUser, err := s.users.GetByID(ctx, leavingReviewerID)
	if err != nil {
		return err
	}

	teamActive, err := s.users.ListActiveUsersInTeam(ctx, leavingUser.TeamName)
	if err != nil {
		return err
	}

	candidates := make([]uuid.UUID, 0)

	for _, userID := range teamActive {
		if userID != leavingReviewerID && userID != pullRequest.AuthorID {
			candidates = append(candidates, userID)
		}
	}

	replacement, err := s.reassignStrategy.PickReplacementReviewers(
		ctx,
		pullRequest,
		leavingReviewerID,
		candidates,
	)
	if err != nil {
		return err
	}

	if replacement == uuid.Nil {
		return apperrors.ErrNoCandidate
	}

	if err = s.pullRequests.AddReviewer(ctx, pullRequestID, replacement); err != nil {
		return err
	}

	if err = s.pullRequests.RemoveReviewer(ctx, pullRequestID, leavingReviewerID); err != nil {
		return err
	}

	return nil
}
