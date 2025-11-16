package mappers

import (
	"fmt"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
)

func ToDomainTeamMember(userDB generated.User) (*models.TeamMember, error) {
	teamMemberDomain, err := models.NewTeamMember(
		UUIDFromPg(userDB.UserID),
		SafeString(userDB.Username),
		userDB.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed create team member domain model: %w", err)
	}

	return teamMemberDomain, nil
}
