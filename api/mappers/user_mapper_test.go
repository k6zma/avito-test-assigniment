package mappers_test

import (
	"testing"

	"github.com/google/uuid"

	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

func TestToDomainUser_Valid(t *testing.T) {
	initValidators(t)

	id := uuid.New()
	dto := dtos.User{
		UserId:   id.String(),
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	user, err := mappers.ToDomainUser(dto)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if user.ID != id {
		t.Errorf("expected ID %v but got %v", id, user.ID)
	}

	if user.Name != dto.Username {
		t.Errorf("expected Name %q but got %q", dto.Username, user.Name)
	}

	if user.TeamName != dto.TeamName {
		t.Errorf("expected TeamName %q but got %q", dto.TeamName, user.TeamName)
	}

	if user.IsActive != dto.IsActive {
		t.Errorf("expected IsActive %v but got %v", dto.IsActive, user.IsActive)
	}
}

func TestToDomainUser_InvalidUUID(t *testing.T) {
	initValidators(t)

	dto := dtos.User{
		UserId:   "bad-uuid",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	if _, err := mappers.ToDomainUser(dto); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func TestToDomainUser_InvalidModel(t *testing.T) {
	initValidators(t)

	dto := dtos.User{
		UserId:   uuid.New().String(),
		Username: "Alice",
		TeamName: "",
		IsActive: true,
	}

	if _, err := mappers.ToDomainUser(dto); err == nil {
		t.Fatalf("expected validation error but got nil")
	}
}

func TestToDTOUser(t *testing.T) {
	initValidators(t)

	id := uuid.New()
	user, err := models.NewUser(id, "Bob", "devops", false)
	if err != nil {
		t.Fatalf("setup user failed: %v", err)
	}

	dto := mappers.ToDTOUser(user)

	if dto.UserId != id.String() {
		t.Errorf("expected UserId %q but got %q", id.String(), dto.UserId)
	}

	if dto.Username != user.Name {
		t.Errorf("expected Username %q but got %q", user.Name, dto.Username)
	}

	if dto.TeamName != user.TeamName {
		t.Errorf("expected TeamName %q but got %q", user.TeamName, dto.TeamName)
	}

	if dto.IsActive != user.IsActive {
		t.Errorf("expected IsActive %v but got %v", user.IsActive, dto.IsActive)
	}
}
