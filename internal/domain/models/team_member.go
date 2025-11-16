package models

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

type TeamMember struct {
	ID       uuid.UUID `validate:"required,uuid"`
	Name     string    `validate:"required,printascii"`
	IsActive bool
}

func NewTeamMember(
	id uuid.UUID,
	name string,
	isActive bool,
) (*TeamMember, error) {
	teamMember := &TeamMember{
		ID:       id,
		Name:     name,
		IsActive: isActive,
	}

	if err := teamMember.Validate(); err != nil {
		return nil, err
	}

	return teamMember, nil
}

func (tm *TeamMember) Validate() error {
	if err := validators.Validator.Struct(tm); err != nil {
		return fmt.Errorf("provided bad team member domain model: %w", err)
	}

	return nil
}
