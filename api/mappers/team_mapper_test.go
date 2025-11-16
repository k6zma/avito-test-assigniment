package mappers_test

import (
	"testing"

	"github.com/google/uuid"

	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

func validTeamMember(t *testing.T, name string, isActive bool) dtos.TeamMember {
	t.Helper()

	return dtos.TeamMember{
		UserId:   uuid.New().String(),
		Username: name,
		IsActive: isActive,
	}
}

func TestToDomainTeam_Valid(t *testing.T) {
	initValidators(t)

	dto := dtos.Team{
		TeamName: "backend",
		Members: []dtos.TeamMember{
			validTeamMember(t, "Alice", true),
			validTeamMember(t, "Bob", true),
		},
	}

	team, err := mappers.ToDomainTeam(dto)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if team.Name != dto.TeamName {
		t.Errorf("expected TeamName %q but got %q", dto.TeamName, team.Name)
	}

	if len(team.Members) != len(dto.Members) {
		t.Errorf("expected %d members but got %d", len(dto.Members), len(team.Members))
	}
}

func TestToDomainTeam_InvalidMember(t *testing.T) {
	initValidators(t)

	dto := dtos.Team{
		TeamName: "backend",
		Members: []dtos.TeamMember{
			{UserId: "bad-uuid", Username: "Alice", IsActive: true},
		},
	}

	if _, err := mappers.ToDomainTeam(dto); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func TestToDomainTeam_InvalidTeam(t *testing.T) {
	initValidators(t)

	dto := dtos.Team{
		TeamName: "",
		Members:  []dtos.TeamMember{validTeamMember(t, "Alice", true)},
	}

	if _, err := mappers.ToDomainTeam(dto); err == nil {
		t.Fatalf("expected validation error but got nil")
	}
}

func TestToDTOTeam(t *testing.T) {
	initValidators(t)

	member1, err := models.NewTeamMember(uuid.New(), "Alice", true)
	if err != nil {
		t.Fatalf("setup member failed: %v", err)
	}

	member2, err := models.NewTeamMember(uuid.New(), "Bob", false)
	if err != nil {
		t.Fatalf("setup member failed: %v", err)
	}

	team, err := models.NewTeam("data", []*models.TeamMember{member1, member2})
	if err != nil {
		t.Fatalf("setup team failed: %v", err)
	}

	dto := mappers.ToDTOTeam(team)

	if dto.TeamName != team.Name {
		t.Errorf("expected TeamName %q but got %q", team.Name, dto.TeamName)
	}

	if len(dto.Members) != len(team.Members) {
		t.Errorf("expected %d members but got %d", len(team.Members), len(dto.Members))
	}
}

