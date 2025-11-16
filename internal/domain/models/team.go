package models

import (
	"fmt"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

type Team struct {
	Name    string        `validate:"required,printascii"`
	Members []*TeamMember `validate:"required,min=1,dive,required"`
}

func NewTeam(
	name string,
	members []*TeamMember,
) (*Team, error) {
	team := &Team{
		Name:    name,
		Members: members,
	}

	if err := team.Validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *Team) Validate() error {
	if err := validators.Validator.Struct(t); err != nil {
		return fmt.Errorf("provided bad team domain model: %w", err)
	}

	return nil
}
