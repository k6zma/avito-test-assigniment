package models_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type teamMemberTestCase struct {
	name     string
	id       uuid.UUID
	member   string
	isActive bool
	wantErr  bool
}

func TestNewTeamMember_Valid(t *testing.T) {
	initValidators(t)

	id := uuid.New()
	member, err := models.NewTeamMember(id, "Alice", true)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if member.ID != id {
		t.Errorf("expected ID %v but got %v", id, member.ID)
	}

	if member.Name != "Alice" {
		t.Errorf("expected Name 'Alice' but got %q", member.Name)
	}

	if !member.IsActive {
		t.Errorf("expected IsActive to be true")
	}
}

func TestNewTeamMember_Invalid(t *testing.T) {
	initValidators(t)

	tests := []teamMemberTestCase{
		{"empty name", uuid.New(), "", true, true},
		{"invalid uuid", uuid.Nil, "Bob", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := models.NewTeamMember(tt.id, tt.member, tt.isActive)
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

func TestTeamMember_Validate(t *testing.T) {
	initValidators(t)

	member := &models.TeamMember{
		ID:       uuid.New(),
		Name:     "Charlie",
		IsActive: true,
	}

	if err := member.Validate(); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}
