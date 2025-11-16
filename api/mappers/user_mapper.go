package mappers

import (
	"fmt"

	"github.com/google/uuid"
	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

func ToDomainUser(dto dtos.User) (*models.User, error) {
	id, err := uuid.Parse(dto.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	user, err := models.NewUser(
		id,
		dto.Username,
		dto.TeamName,
		dto.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid user dto: %w", err)
	}

	return user, nil
}

func ToDTOUser(user *models.User) dtos.User {
	return dtos.User{
		UserId:   user.ID.String(),
		Username: user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}
