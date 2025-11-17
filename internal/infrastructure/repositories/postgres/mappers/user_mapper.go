package mappers

import (
	"errors"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
)

func ToDomainUser(userDB generated.User) (*models.User, error) {
	userDomain, err := models.NewUser(
		UUIDFromPg(userDB.UserID),
		SafeString(userDB.Username),
		userDB.TeamName,
		userDB.IsActive,
	)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, err)
	}

	return userDomain, nil
}

func ToDBUser(userDomain *models.User) generated.UpsertUserParams {
	return generated.UpsertUserParams{
		UserID:   UUIDToPg(userDomain.ID),
		Username: StringPtr(userDomain.Name),
		TeamName: userDomain.TeamName,
		IsActive: userDomain.IsActive,
	}
}
