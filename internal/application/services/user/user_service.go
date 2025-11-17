package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/pkg/logger"
)

type UserService struct {
	users  repositories.UserRepository
	teams  repositories.TeamRepository
	logger *slog.Logger
}

func NewUserService(
	userRepo repositories.UserRepository,
	teamRepo repositories.TeamRepository,
	logger *slog.Logger,
) *UserService {
	const userEntity = "user"

	return &UserService{
		users:  userRepo,
		teams:  teamRepo,
		logger: logger.WithGroup(userEntity),
	}
}

func (s *UserService) UpsertUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	teamName string,
	isActive bool,
) (user *models.User, err error) {
	const upsertUserMethod = "UpsertUser"

	log := s.logger.WithGroup(upsertUserMethod)

	operationLog := logger.StartOperation(
		log,
		"upserting user",
		slog.String("user_id", userID.String()),
		slog.String("team_name", teamName),
		slog.String("username", name),
		slog.Bool("is_active", isActive),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("user_id", userID.String()),
		slog.String("team_name", teamName),
		slog.String("username", name),
		slog.Bool("is_active", isActive),
	)

	if existing, err := s.users.GetByID(ctx, userID); err == nil {
		if existing.TeamName != teamName {
			return nil, errors.Join(
				apperrors.ErrValidation,
				fmt.Errorf("user %s already belongs to team %s", userID, existing.TeamName),
			)
		}
	}

	user, err = models.NewUser(userID, name, teamName, isActive)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, err)
	}

	if err = s.teams.Create(ctx, teamName); err != nil {
		if !errors.Is(err, apperrors.ErrTeamExists) {
			return nil, err
		}
	}

	if err = s.users.Upsert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(
	ctx context.Context,
	userID uuid.UUID,
) (user *models.User, err error) {
	const getUserMethod = "GetUser"

	log := s.logger.WithGroup(getUserMethod)

	operationLog := logger.StartOperation(
		log,
		"getting user",
		slog.String("user_id", userID.String()),
	)
	defer operationLog.FinishOperation(&err, slog.String("user_id", userID.String()))

	user, err = s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed get user by id: %w", err)
	}

	return user, nil
}

func (s *UserService) SetUserActive(
	ctx context.Context,
	userID uuid.UUID,
	isActive bool,
) (user *models.User, err error) {
	const setUserActiveMethod = "SetUserActive"

	log := s.logger.WithGroup(setUserActiveMethod)

	operationLog := logger.StartOperation(
		log,
		"setting user activity",
		slog.String("user_id", userID.String()),
		slog.Bool("is_active", isActive),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("user_id", userID.String()),
		slog.Bool("is_active", isActive),
	)

	user, err = s.users.SetActive(ctx, userID, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed update user active flag: %w", err)
	}

	return user, nil
}

func (s *UserService) ListActiveUsersInTeam(
	ctx context.Context,
	teamName string,
) (users []uuid.UUID, err error) {
	const listActiveUserInTeamMethod = "ListActiveUsersInTeam"

	log := s.logger.WithGroup(listActiveUserInTeamMethod)

	operationLog := logger.StartOperation(
		log,
		"setting user activity",
		slog.String("team_name", teamName),
	)
	defer operationLog.FinishOperation(
		&err,
		slog.String("team_name", teamName),
	)

	userIDS, err := s.users.ListActiveUsersInTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed list active users in team: %w", err)
	}

	return userIDS, nil
}
