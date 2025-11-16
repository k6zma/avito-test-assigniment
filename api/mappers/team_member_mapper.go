package mappers

import (
	"fmt"

	"github.com/google/uuid"
	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

func ToDomainTeamMember(dto dtos.TeamMember) (*models.TeamMember, error) {
	id, err := uuid.Parse(dto.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	teamMember, err := models.NewTeamMember(id, dto.Username, dto.IsActive)
	if err != nil {
		return nil, fmt.Errorf("invalid team member dto: %w", err)
	}

	return teamMember, nil
}

func ToDTOTeamMember(teamMembers *models.TeamMember) dtos.TeamMember {
	return dtos.TeamMember{
		UserId:   teamMembers.ID.String(),
		Username: teamMembers.Name,
		IsActive: teamMembers.IsActive,
	}
}
