package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/valueobjects"
	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

type PullRequest struct {
	ID                uuid.UUID                      `validate:"required,uuid"`
	Name              string                         `validate:"required,printascii"`
	AuthorID          uuid.UUID                      `validate:"required,uuid"`
	Status            valueobjects.PullRequestStatus `validate:"required,oneof=OPEN MERGED"`
	AssignedReviewers []uuid.UUID                    `validate:"dive,required,uuid"`
	CreatedAt         time.Time                      `validate:"required"`
	MergedAt          *time.Time                     `validate:"omitempty"`
}

func NewPullRequest(
	id uuid.UUID,
	name string,
	authorID uuid.UUID,
	status valueobjects.PullRequestStatus,
	assignedReviewers []uuid.UUID,
) (*PullRequest, error) {
	pr := &PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            status,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         time.Now().UTC(),
	}

	if err := pr.Validate(); err != nil {
		return nil, err
	}

	return pr, nil
}

func (pr *PullRequest) Validate() error {
	if err := validators.Validator.Struct(pr); err != nil {
		return fmt.Errorf("provided bad pull request domain model: %w", err)
	}

	return nil
}
