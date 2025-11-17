package mappers_test

import (
	"testing"

	"github.com/google/uuid"

	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

func TestToDomainTeamMember_Valid(t *testing.T) {
	initValidators(t)

	id := uuid.New()
	dto := dtos.TeamMember{
		UserId:   id.String(),
		Username: "Alice",
		IsActive: true,
	}

	member, err := mappers.ToDomainTeamMember(dto)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if member.ID != id {
		t.Errorf("expected ID %v but got %v", id, member.ID)
	}

	if member.Name != dto.Username {
		t.Errorf("expected Name %q but got %q", dto.Username, member.Name)
	}

	if member.IsActive != dto.IsActive {
		t.Errorf("expected IsActive %v but got %v", dto.IsActive, member.IsActive)
	}
}

func TestToDomainTeamMember_InvalidUUID(t *testing.T) {
	initValidators(t)

	dto := dtos.TeamMember{
		UserId:   "not-a-uuid",
		Username: "Alice",
		IsActive: true,
	}

	if _, err := mappers.ToDomainTeamMember(dto); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func TestToDomainTeamMember_InvalidModel(t *testing.T) {
	initValidators(t)

	dto := dtos.TeamMember{
		UserId:   uuid.New().String(),
		Username: "",
		IsActive: true,
	}

	if _, err := mappers.ToDomainTeamMember(dto); err == nil {
		t.Fatalf("expected validation error but got nil")
	}
}

func TestToDTOTeamMember(t *testing.T) {
	initValidators(t)

	id := uuid.New()

	member, err := models.NewTeamMember(id, "Bob", false)
	if err != nil {
		t.Fatalf("setup team member failed: %v", err)
	}

	dto := mappers.ToDTOTeamMember(member)

	if dto.UserId != id.String() {
		t.Errorf("expected UserId %q but got %q", id.String(), dto.UserId)
	}

	if dto.Username != member.Name {
		t.Errorf("expected Username %q but got %q", member.Name, dto.Username)
	}

	if dto.IsActive != member.IsActive {
		t.Errorf("expected IsActive %v but got %v", member.IsActive, dto.IsActive)
	}
}
