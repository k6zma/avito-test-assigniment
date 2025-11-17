package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/mappers"
	"github.com/k6zma/avito-test-assigniment/pkg/logger"
)

type UserRepository struct {
	queries *generated.Queries
	logger  *slog.Logger
}

func NewUserRepository(
	q *generated.Queries,
	logger *slog.Logger,
) repositories.UserRepository {
	const userEntity = "user"

	return &UserRepository{
		queries: q,
		logger:  logger.WithGroup(userEntity).WithGroup(postgresDBType),
	}
}

func (r *UserRepository) Upsert(ctx context.Context, user *models.User) (err error) {
	const upsertMethod = "Upsert"

	params := mappers.ToDBUser(user)

	log := r.logger.WithGroup(upsertMethod)

	operationLog := logger.StartOperation(
		log,
		"putting data in storage",
		slog.String("user_id", user.ID.String()),
		slog.String("user_name", user.Name),
		slog.String("team_name", user.TeamName),
		slog.Bool("is_active", user.IsActive),
	)

	defer func() { operationLog.FinishOperation(&err) }()

	if err = r.queries.UpsertUser(ctx, params); err != nil {
		return fmt.Errorf("failed upsert user in postgres storage: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	userID uuid.UUID,
) (user *models.User, err error) {
	const getByIDMethod = "GetByID"

	log := r.logger.WithGroup(getByIDMethod)

	operationLog := logger.StartOperation(
		log,
		"putting data in storage",
		slog.String("user_id", userID.String()),
	)

	defer func() {
		var extra []any

		if user != nil {
			extra = append(
				extra,
				slog.String("user_id", user.ID.String()),
				slog.String("user_name", user.Name),
				slog.String("team_name", user.TeamName),
				slog.Bool("is_active", user.IsActive),
			)
		}

		operationLog.FinishOperation(&err, extra...)
	}()

	dbUser, err := r.queries.GetUser(ctx, mappers.UUIDToPg(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrUserNotFound, err)
		}

		return nil, fmt.Errorf("failed get user from postgres storage: %w", err)
	}

	user, err = mappers.ToDomainUser(dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed convert user to domain model: %w", err)
	}

	return user, nil
}

func (r *UserRepository) ListActiveUsersInTeam(
	ctx context.Context,
	teamName string,
) (userIDs []uuid.UUID, err error) {
	const listActiveUsersMethod = "ListActiveUsersInTeam"

	log := r.logger.WithGroup(listActiveUsersMethod)

	operationLog := logger.StartOperation(
		log,
		"getting active users",
		slog.String("team_name", teamName),
	)

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.String("team_name", teamName),
			slog.Int("users_count", len(userIDs)),
		)
	}()

	userIDS, err := r.queries.ListActiveUsersInTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed get list active users in team from postgres storage: %w",
			err,
		)
	}

	userIDs = make([]uuid.UUID, 0, len(userIDS))

	for _, userID := range userIDS {
		userIDs = append(userIDs, mappers.UUIDFromPg(userID))
	}

	return userIDs, nil
}

func (r *UserRepository) ListUsersInTeam(
	ctx context.Context,
	teamName string,
) (userIDs []uuid.UUID, err error) {
	const listUsersMethod = "ListUsersInTeam"

	log := r.logger.WithGroup(listUsersMethod)

	operationLog := logger.StartOperation(
		log,
		"getting users",
		slog.String("team_name", teamName),
	)

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.String("team_name", teamName),
			slog.Int("users_count", len(userIDs)),
		)
	}()

	userIDS, err := r.queries.ListUsersInTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed get list users in team from postgres storage: %w",
			err,
		)
	}

	userIDs = make([]uuid.UUID, 0, len(userIDS))

	for _, userID := range userIDS {
		userIDs = append(userIDs, mappers.UUIDFromPg(userID))
	}

	return userIDs, nil
}

func (r *UserRepository) SetActive(
	ctx context.Context,
	userID uuid.UUID,
	isActive bool,
) (user *models.User, err error) {
	const setActiveMethod = "SetActive"

	log := r.logger.WithGroup(setActiveMethod)

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

	dbUser, err := r.queries.SetUserActive(ctx, generated.SetUserActiveParams{
		UserID:   mappers.UUIDToPg(userID),
		IsActive: isActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrUserNotFound, err)
		}

		return nil, fmt.Errorf("failed set user active in postgres storage: %w", err)
	}

	user, err = mappers.ToDomainUser(dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed convert user to domain model: %w", err)
	}

	return user, nil
}
