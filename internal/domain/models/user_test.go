package models_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type userTestCase struct {
	name     string
	id       uuid.UUID
	username string
	teamName string
	isActive bool
	wantErr  bool
}

func TestNewUser_Valid(t *testing.T) {
	initValidators(t)

	id := uuid.New()

	user, err := models.NewUser(id, "Alice", "backend", true)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if user.ID != id {
		t.Errorf("expected ID %v but got %v", id, user.ID)
	}

	if user.Name != "Alice" {
		t.Errorf("expected Name 'Alice' but got %q", user.Name)
	}

	if user.TeamName != "backend" {
		t.Errorf("expected TeamName 'backend' but got %q", user.TeamName)
	}
}

func TestNewUser_Invalid(t *testing.T) {
	initValidators(t)

	tests := []userTestCase{
		{"empty name", uuid.New(), "", "backend", true, true},
		{"empty team name", uuid.New(), "Alice", "", true, true},
		{"invalid uuid", uuid.Nil, "Alice", "backend", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := models.NewUser(tt.id, tt.username, tt.teamName, tt.isActive)
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

func TestUser_Validate(t *testing.T) {
	initValidators(t)

	user := &models.User{
		ID:       uuid.New(),
		Name:     "Bob",
		TeamName: "devops",
		IsActive: true,
	}

	if err := user.Validate(); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}
