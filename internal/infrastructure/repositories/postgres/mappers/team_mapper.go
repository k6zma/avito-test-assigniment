package mappers

import (
	"errors"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
)

func ToDomainTeam(teamDB generated.Team, members []*models.TeamMember) (*models.Team, error) {
	teamDomain, err := models.NewTeam(teamDB.TeamName, members)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, err)
	}

	return teamDomain, nil
}
