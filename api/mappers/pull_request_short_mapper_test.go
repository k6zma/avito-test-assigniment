package mappers_test

import (
	"testing"

	"github.com/google/uuid"

	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
)

func TestToDTOPRShort(t *testing.T) {
	initValidators(t)

	id := uuid.New()
	authorID := uuid.New()

	pr, err := models.NewPullRequest(id, "Title", authorID, valueobjects.OpenPullRequest, nil)
	if err != nil {
		t.Fatalf("setup pull request failed: %v", err)
	}

	dto := mappers.ToDTOPRShort(pr)

	if dto.PullRequestId != pr.ID.String() {
		t.Errorf("expected PullRequestId %q but got %q", pr.ID.String(), dto.PullRequestId)
	}

	if dto.PullRequestName != pr.Name {
		t.Errorf("expected PullRequestName %q but got %q", pr.Name, dto.PullRequestName)
	}

	if dto.AuthorId != pr.AuthorID.String() {
		t.Errorf("expected AuthorId %q but got %q", pr.AuthorID.String(), dto.AuthorId)
	}

	if dto.Status != dtos.PullRequestShortStatus(pr.Status) {
		t.Errorf("expected Status %q but got %q", pr.Status, dto.Status)
	}
}

