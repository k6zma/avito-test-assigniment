package mappers

import (
	"fmt"

	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

func ToDomainTeam(dto dtos.Team) (*models.Team, error) {
	members := make([]*models.TeamMember, 0, len(dto.Members))

	for _, memberDTO := range dto.Members {
		teamMember, err := ToDomainTeamMember(memberDTO)
		if err != nil {
			return nil, fmt.Errorf("invalid member: %w", err)
		}

		members = append(members, teamMember)
	}

	team, err := models.NewTeam(dto.TeamName, members)
	if err != nil {
		return nil, fmt.Errorf("invalid team dto: %w", err)
	}

	return team, nil
}

func ToDTOTeam(team *models.Team) dtos.Team {
	dtoMembers := make([]dtos.TeamMember, 0, len(team.Members))

	for _, member := range team.Members {
		dtoMembers = append(dtoMembers, ToDTOTeamMember(member))
	}

	return dtos.Team{
		TeamName: team.Name,
		Members:  dtoMembers,
	}
}
