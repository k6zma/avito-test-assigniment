package models

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

type User struct {
	ID       uuid.UUID `validate:"required,uuid"`
	Name     string    `validate:"required,printascii"`
	TeamName string    `validate:"required,printascii"`
	IsActive bool
}

func NewUser(
	id uuid.UUID,
	name string,
	teamName string,
	isActive bool,
) (*User, error) {
	user := &User{
		ID:       id,
		Name:     name,
		TeamName: teamName,
		IsActive: isActive,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Validate() error {
	if err := validators.Validator.Struct(u); err != nil {
		return fmt.Errorf("provided bad user domain model: %w", err)
	}

	return nil
}
