package pullrequest

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest/reviewers/initialize"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest/reviewers/reassign"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
)

type PullRequestService struct {
	pullRequests     repositories.PullRequestRepository
	users            repositories.UserRepository
	initStrategy     initialize.ReviewerInitializeStrategy
	reassignStrategy reassign.ReviewerReassignStrategy
}

func NewPullRequestService(
	prRepo repositories.PullRequestRepository,
	userRepo repositories.UserRepository,
	initStrategy initialize.ReviewerInitializeStrategy,
	reassignStrategy reassign.ReviewerReassignStrategy,
) *PullRequestService {
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
	}
}

func (s *PullRequestService) CreatePullRequest(
	ctx context.Context,
	id uuid.UUID,
	name string,
	authorID uuid.UUID,
) (*models.PullRequest, error) {
	author, err := s.users.GetByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed get author user: %w", err)
	}
	if !author.IsActive {
		return nil, fmt.Errorf("author %s is not active", authorID)
	}

	teamActive, err := s.users.ListActiveUsersInTeam(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed list active users in team: %w", err)
	}

	candidates := make([]uuid.UUID, 0, len(teamActive))

	for _, uid := range teamActive {
		if uid != authorID {
			candidates = append(candidates, uid)
		}
	}

	pr, err := models.NewPullRequest(
		id,
		name,
		authorID,
		valueobjects.OpenPullRequest,
		[]uuid.UUID{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed create pull request domain model: %w", err)
	}

	assigned, err := s.initStrategy.PickInitialReviewers(ctx, pr, candidates)
	if err != nil {
		return nil, fmt.Errorf("failed pick initial reviewers: %w", err)
	}

	pr.AssignedReviewers = assigned

	if err = s.pullRequests.Create(ctx, pr); err != nil {
		return nil, fmt.Errorf("failed create pull request: %w", err)
	}

	return pr, nil
}

func (s *PullRequestService) GetPullRequest(
	ctx context.Context,
	id uuid.UUID,
) (*models.PullRequest, error) {
	pr, err := s.pullRequests.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed get pull request by id: %w", err)
	}

	return pr, nil
}

func (s *PullRequestService) MergePullRequest(
	ctx context.Context,
	id uuid.UUID,
) (*models.PullRequest, error) {
	pr, err := s.pullRequests.Merge(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed merge pull request: %w", err)
	}

	return pr, nil
}

func (s *PullRequestService) ListPullRequestsAssignedToUser(
	ctx context.Context,
	reviewerID uuid.UUID,
) ([]*models.PullRequest, error) {
	prs, err := s.pullRequests.ListAssignedToUser(ctx, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed list pull requests assigned to user: %w", err)
	}

	return prs, nil
}

func (s *PullRequestService) AddReviewer(
	ctx context.Context,
	pullRequestID, reviewerID uuid.UUID,
) error {
	pr, err := s.pullRequests.GetByID(ctx, pullRequestID)
	if err != nil {
		return fmt.Errorf("failed get pull request: %w", err)
	}

	if pr.Status != valueobjects.OpenPullRequest {
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
) error {
	pr, err := s.pullRequests.GetByID(ctx, pullRequestID)
	if err != nil {
		return fmt.Errorf("failed get pull request: %w", err)
	}

	if pr.Status != valueobjects.OpenPullRequest {
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
) ([]uuid.UUID, error) {
	revs, err := s.pullRequests.ListReviewers(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed list reviewers for pull request: %w", err)
	}

	return revs, nil
}

func (s *PullRequestService) CountReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (int64, error) {
	count, err := s.pullRequests.CountReviewers(ctx, pullRequestID)
	if err != nil {
		return 0, fmt.Errorf("failed count reviewers for pull request: %w", err)
	}

	return count, nil
}

func (s *PullRequestService) ReassignReviewer(
	ctx context.Context,
	pullRequestID, leavingReviewerID uuid.UUID,
) error {
	pr, err := s.pullRequests.GetByID(ctx, pullRequestID)
	if err != nil {
		return fmt.Errorf("failed get pull request: %w", err)
	}

	if pr.Status != valueobjects.OpenPullRequest {
		return fmt.Errorf("cannot reassign reviewers for non-open pull request")
	}

	found := false

	for _, rid := range pr.AssignedReviewers {
		if rid == leavingReviewerID {
			found = true

			break
		}
	}

	if !found {
		return fmt.Errorf("user %s is not a reviewer of PR %s", leavingReviewerID, pullRequestID)
	}

	leavingUser, err := s.users.GetByID(ctx, leavingReviewerID)
	if err != nil {
		return fmt.Errorf("failed get leaving reviewer: %w", err)
	}

	teamActive, err := s.users.ListActiveUsersInTeam(ctx, leavingUser.TeamName)
	if err != nil {
		return fmt.Errorf("failed list active users: %w", err)
	}

	candidates := make([]uuid.UUID, 0)

	for _, uid := range teamActive {
		if uid != leavingReviewerID && uid != pr.AuthorID {
			candidates = append(candidates, uid)
		}
	}

	replacement, err := s.reassignStrategy.PickReplacementReviewers(
		ctx,
		pr,
		leavingReviewerID,
		candidates,
	)
	if err != nil {
		return fmt.Errorf("failed pick replacement reviewer: %w", err)
	}

	if replacement == uuid.Nil {
		return fmt.Errorf("no available candidates to reassign reviewer")
	}

	if err = s.pullRequests.AddReviewer(ctx, pullRequestID, replacement); err != nil {
		return fmt.Errorf("failed add replacement reviewer: %w", err)
	}

	if err = s.pullRequests.RemoveReviewer(ctx, pullRequestID, leavingReviewerID); err != nil {
		return fmt.Errorf("failed remove old reviewer: %w", err)
	}

	return nil
}
