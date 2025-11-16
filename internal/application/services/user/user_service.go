package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
)

type UserService struct {
	users repositories.UserRepository
	teams repositories.TeamRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	teamRepo repositories.TeamRepository,
) *UserService {
	return &UserService{
		users: userRepo,
		teams: teamRepo,
	}
}

func (s *UserService) UpsertUser(
	ctx context.Context,
	id uuid.UUID,
	name string,
	teamName string,
	isActive bool,
) (*models.User, error) {
	user, err := models.NewUser(id, name, teamName, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed create user domain model: %w", err)
	}

	if err = s.teams.Create(ctx, teamName); err != nil {
		return nil, fmt.Errorf("failed ensure team exists: %w", err)
	}

	if err = s.users.Upsert(ctx, user); err != nil {
		return nil, fmt.Errorf("failed upsert user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed get user by id: %w", err)
	}

	return user, nil
}

func (s *UserService) SetUserActive(
	ctx context.Context,
	id uuid.UUID,
	isActive bool,
) (*models.User, error) {
	user, err := s.users.SetActive(ctx, id, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed update user active flag: %w", err)
	}

	return user, nil
}

func (s *UserService) ListActiveUsersInTeam(
	ctx context.Context,
	teamName string,
) ([]uuid.UUID, error) {
	ids, err := s.users.ListActiveUsersInTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed list active users in team: %w", err)
	}

	return ids, nil
}
