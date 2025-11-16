package models_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type teamTestCase struct {
	name    string
	team    string
	members []*models.TeamMember
	wantErr bool
}

func validMember(t *testing.T, name string) *models.TeamMember {
	t.Helper()

	member, err := models.NewTeamMember(uuid.New(), name, true)
	if err != nil {
		t.Fatalf("failed to create test team member: %v", err)
	}

	return member
}

func TestNewTeam_Valid(t *testing.T) {
	initValidators(t)

	members := []*models.TeamMember{
		validMember(t, "Alice"),
		validMember(t, "Bob"),
	}

	team, err := models.NewTeam("backend", members)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if team.Name != "backend" {
		t.Errorf("expected team name 'backend' but got %q", team.Name)
	}

	if len(team.Members) != len(members) {
		t.Errorf("expected %d members but got %d", len(members), len(team.Members))
	}
}

func TestNewTeam_Invalid(t *testing.T) {
	initValidators(t)

	valid := validMember(t, "Alice")

	tests := []teamTestCase{
		{"empty name", "", []*models.TeamMember{valid}, true},
		{"nil members slice", "devops", nil, true},
		{"empty members", "qa", []*models.TeamMember{}, true},
		{"nil member entry", "analytics", []*models.TeamMember{nil}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := models.NewTeam(tt.team, tt.members)
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

func TestTeam_Validate(t *testing.T) {
	initValidators(t)

	team := &models.Team{
		Name: "mobile",
		Members: []*models.TeamMember{
			validMember(t, "Alice"),
		},
	}

	if err := team.Validate(); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}
