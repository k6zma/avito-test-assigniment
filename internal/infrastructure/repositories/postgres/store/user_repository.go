package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/mappers"
)

type UserRepository struct {
	queries *generated.Queries
}

func NewUserRepository(q *generated.Queries) repositories.UserRepository {
	return &UserRepository{queries: q}
}

func (r *UserRepository) Upsert(ctx context.Context, user *models.User) error {
	params := mappers.ToDBUser(user)

	err := r.queries.UpsertUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed upsert user in postgres storage: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	dbUser, err := r.queries.GetUser(ctx, mappers.UUIDToPg(id))
	if err != nil {
		return nil, fmt.Errorf("failed get user from postgres storage: %w", err)
	}

	domainUser, err := mappers.ToDomainUser(dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed convert user to domain model: %w", err)
	}

	return domainUser, nil
}

func (r *UserRepository) ListActiveUsersInTeam(
	ctx context.Context,
	teamName string,
) ([]uuid.UUID, error) {
	ids, err := r.queries.ListActiveUsersInTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed get list active users in team from postgres storage: %w",
			err,
		)
	}

	result := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		result = append(result, mappers.UUIDFromPg(id))
	}

	return result, nil
}

func (r *UserRepository) SetActive(
	ctx context.Context,
	id uuid.UUID,
	isActive bool,
) (*models.User, error) {
	dbUser, err := r.queries.SetUserActive(ctx, generated.SetUserActiveParams{
		UserID:   mappers.UUIDToPg(id),
		IsActive: isActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed set user active in postgress storage: %w", err)
	}

	domainUser, err := mappers.ToDomainUser(dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed convert user to domain model: %w", err)
	}

	return domainUser, nil
}
