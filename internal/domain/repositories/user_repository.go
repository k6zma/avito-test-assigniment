package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

type UserRepository interface {
	Upsert(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListActiveUsersInTeam(ctx context.Context, teamName string) ([]uuid.UUID, error)
	ListUsersInTeam(ctx context.Context, teamName string) ([]uuid.UUID, error)
	SetActive(ctx context.Context, id uuid.UUID, isActive bool) (*models.User, error)
}
