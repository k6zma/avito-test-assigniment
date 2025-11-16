package mappers_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
)

type pullRequestDTOTestCase struct {
	name    string
	dto     dtos.PullRequest
	wantErr bool
}

func TestToDomainPullRequest_Valid(t *testing.T) {
	initValidators(t)

	id := uuid.New()
	authorID := uuid.New()
	reviewer1 := uuid.New()
	reviewer2 := uuid.New()
	createdAt := time.Now().UTC()
	mergedAt := createdAt.Add(time.Hour).UTC()

	dto := dtos.PullRequest{
		PullRequestId:     id.String(),
		PullRequestName:   "Add feature",
		AuthorId:          authorID.String(),
		Status:            dtos.PullRequestStatusOPEN,
		AssignedReviewers: []string{reviewer1.String(), reviewer2.String()},
		CreatedAt:         &createdAt,
		MergedAt:          &mergedAt,
	}

	pr, err := mappers.ToDomainPullRequest(dto)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if pr.ID != id {
		t.Errorf("expected ID %v but got %v", id, pr.ID)
	}

	if pr.AuthorID != authorID {
		t.Errorf("expected AuthorID %v but got %v", authorID, pr.AuthorID)
	}

	if pr.Status != valueobjects.OpenPullRequest {
		t.Errorf("expected Status %q but got %q", valueobjects.OpenPullRequest, pr.Status)
	}

	if len(pr.AssignedReviewers) != 2 {
		t.Errorf("expected 2 reviewers but got %d", len(pr.AssignedReviewers))
	}

	if pr.MergedAt == nil || !pr.MergedAt.Equal(mergedAt) {
		t.Errorf("expected MergedAt %v but got %v", mergedAt, pr.MergedAt)
	}

	if !pr.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v but got %v", createdAt, pr.CreatedAt)
	}
}

func TestToDomainPullRequest_InvalidCases(t *testing.T) {
	initValidators(t)

	baseID := uuid.New()
	authorID := uuid.New()
	createdAt := time.Now().UTC()

	tests := []pullRequestDTOTestCase{
		{
			name: "invalid status",
			dto: dtos.PullRequest{
				PullRequestId:     baseID.String(),
				PullRequestName:   "Add feature",
				AuthorId:          authorID.String(),
				Status:            "INVALID",
				AssignedReviewers: []string{},
				CreatedAt:         &createdAt,
			},
			wantErr: true,
		},
		{
			name: "invalid pull request id",
			dto: dtos.PullRequest{
				PullRequestId:     "bad-uuid",
				PullRequestName:   "Add feature",
				AuthorId:          authorID.String(),
				Status:            dtos.PullRequestStatusOPEN,
				AssignedReviewers: []string{},
				CreatedAt:         &createdAt,
			},
			wantErr: true,
		},
		{
			name: "invalid author id",
			dto: dtos.PullRequest{
				PullRequestId:     baseID.String(),
				PullRequestName:   "Add feature",
				AuthorId:          "bad-uuid",
				Status:            dtos.PullRequestStatusOPEN,
				AssignedReviewers: []string{},
				CreatedAt:         &createdAt,
			},
			wantErr: true,
		},
		{
			name: "invalid reviewer id",
			dto: dtos.PullRequest{
				PullRequestId:     baseID.String(),
				PullRequestName:   "Add feature",
				AuthorId:          authorID.String(),
				Status:            dtos.PullRequestStatusOPEN,
				AssignedReviewers: []string{"bad-uuid"},
				CreatedAt:         &createdAt,
			},
			wantErr: true,
		},
		{
			name: "missing created at",
			dto: dtos.PullRequest{
				PullRequestId:     baseID.String(),
				PullRequestName:   "Add feature",
				AuthorId:          authorID.String(),
				Status:            dtos.PullRequestStatusOPEN,
				AssignedReviewers: []string{},
				CreatedAt:         nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mappers.ToDomainPullRequest(tt.dto)
			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"unexpected error state: gotErr=%v, wantErr=%v, err=%v",
					err != nil,
					tt.wantErr,
					err,
				)
			}
		})
	}
}

func TestToDTOPullRequest(t *testing.T) {
	initValidators(t)

	id := uuid.New()
	authorID := uuid.New()
	reviewer := uuid.New()

	pr, err := models.NewPullRequest(
		id,
		"Refactor",
		authorID,
		valueobjects.OpenPullRequest,
		[]uuid.UUID{reviewer},
	)
	if err != nil {
		t.Fatalf("setup pull request failed: %v", err)
	}

	now := time.Now().UTC()
	pr.MergedAt = &now

	dto := mappers.ToDTOPullRequest(pr)

	if dto.PullRequestId != pr.ID.String() {
		t.Errorf("expected PullRequestId %q but got %q", pr.ID.String(), dto.PullRequestId)
	}

	if dto.PullRequestName != pr.Name {
		t.Errorf("expected PullRequestName %q but got %q", pr.Name, dto.PullRequestName)
	}

	if dto.AuthorId != pr.AuthorID.String() {
		t.Errorf("expected AuthorId %q but got %q", pr.AuthorID.String(), dto.AuthorId)
	}

	if len(dto.AssignedReviewers) != len(pr.AssignedReviewers) {
		t.Errorf(
			"expected %d reviewers but got %d",
			len(pr.AssignedReviewers),
			len(dto.AssignedReviewers),
		)
	}

	if dto.CreatedAt == nil || !dto.CreatedAt.Equal(pr.CreatedAt) {
		t.Errorf("expected CreatedAt %v but got %v", pr.CreatedAt, dto.CreatedAt)
	}

	if dto.MergedAt == nil || !dto.MergedAt.Equal(*pr.MergedAt) {
		t.Errorf("expected MergedAt %v but got %v", pr.MergedAt, dto.MergedAt)
	}
}
