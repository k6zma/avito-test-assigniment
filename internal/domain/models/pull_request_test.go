package models_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
)

type pullRequestTestCase struct {
	name      string
	id        uuid.UUID
	title     string
	authorID  uuid.UUID
	status    valueobjects.PullRequestStatus
	reviewers []uuid.UUID
	wantErr   bool
}

func TestNewPullRequest_Valid(t *testing.T) {
	initValidators(t)

	reviewers := []uuid.UUID{uuid.New(), uuid.New()}

	pr, err := models.NewPullRequest(
		uuid.New(),
		"Add feature",
		uuid.New(),
		valueobjects.OpenPullRequest,
		reviewers,
	)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if pr.Name != "Add feature" {
		t.Errorf("expected Name 'Add feature' but got %q", pr.Name)
	}

	if pr.Status != valueobjects.OpenPullRequest {
		t.Errorf("expected Status %q but got %q", valueobjects.OpenPullRequest, pr.Status)
	}

	if len(pr.AssignedReviewers) != len(reviewers) {
		t.Errorf("expected %d reviewers but got %d", len(reviewers), len(pr.AssignedReviewers))
	}

	if pr.MergedAt != nil {
		t.Errorf("expected MergedAt to be nil for new PR")
	}

	if pr.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set")
	}
}

func TestNewPullRequest_ValidWithoutReviewers(t *testing.T) {
	initValidators(t)

	pr, err := models.NewPullRequest(
		uuid.New(),
		"No reviewers",
		uuid.New(),
		valueobjects.OpenPullRequest,
		nil,
	)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if len(pr.AssignedReviewers) != 0 {
		t.Errorf("expected no reviewers but got %d", len(pr.AssignedReviewers))
	}
}

func TestNewPullRequest_Invalid(t *testing.T) {
	initValidators(t)

	tests := []pullRequestTestCase{
		{
			name:      "empty title",
			id:        uuid.New(),
			title:     "",
			authorID:  uuid.New(),
			status:    valueobjects.OpenPullRequest,
			reviewers: []uuid.UUID{uuid.New()},
			wantErr:   true,
		},
		{
			name:      "invalid status",
			id:        uuid.New(),
			title:     "Add feature",
			authorID:  uuid.New(),
			status:    valueobjects.PullRequestStatus("INVALID"),
			reviewers: []uuid.UUID{uuid.New()},
			wantErr:   true,
		},
		{
			name:      "invalid reviewer uuid",
			id:        uuid.New(),
			title:     "Add feature",
			authorID:  uuid.New(),
			status:    valueobjects.OpenPullRequest,
			reviewers: []uuid.UUID{uuid.Nil},
			wantErr:   true,
		},
		{
			name:      "invalid author uuid",
			id:        uuid.New(),
			title:     "Add feature",
			authorID:  uuid.Nil,
			status:    valueobjects.OpenPullRequest,
			reviewers: []uuid.UUID{uuid.New()},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := models.NewPullRequest(tt.id, tt.title, tt.authorID, tt.status, tt.reviewers)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"unexpected error state: gotErr=%v, wantErr=%v, err=%v",
					err != nil,
					tt.wantErr,
					err,
				)
			}
		})
	}
}

func TestPullRequest_Validate(t *testing.T) {
	initValidators(t)

	now := time.Now().UTC()

	pr := &models.PullRequest{
		ID:                uuid.New(),
		Name:              "Merge changes",
		AuthorID:          uuid.New(),
		Status:            valueobjects.MergedPullRequest,
		AssignedReviewers: []uuid.UUID{uuid.New()},
		CreatedAt:         now,
		MergedAt:          &now,
	}

	if err := pr.Validate(); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func TestPullRequest_Validate_InvalidCreatedAt(t *testing.T) {
	initValidators(t)

	pr := &models.PullRequest{
		ID:                uuid.New(),
		Name:              "Invalid created at",
		AuthorID:          uuid.New(),
		Status:            valueobjects.OpenPullRequest,
		AssignedReviewers: []uuid.UUID{},
		CreatedAt:         time.Time{},
		MergedAt:          nil,
	}

	if err := pr.Validate(); err == nil {
		t.Fatalf("expected validation error but got nil")
	}
}
