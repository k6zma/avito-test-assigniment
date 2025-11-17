package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
	"github.com/k6zma/avito-test-assigniment/pkg/logger"
)

type TeamRepository struct {
	queries *generated.Queries
	logger  *slog.Logger
}

func NewTeamRepository(
	q *generated.Queries,
	logger *slog.Logger,
) repositories.TeamRepository {
	const teamEntity = "team"

	return &TeamRepository{
		queries: q,
		logger:  logger.WithGroup(teamEntity).WithGroup(postgresDBType),
	}
}

func (r *TeamRepository) Create(ctx context.Context, name string) (err error) {
	const createMethod = "Create"

	log := r.logger.WithGroup(createMethod)

	operationLog := logger.StartOperation(
		log,
		"creating team in store",
		slog.String("team_name", name),
	)
	defer operationLog.FinishOperation(&err, slog.String("team_name", name))

	err = r.queries.CreateTeam(ctx, name)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}

		return fmt.Errorf("failed create team in postgres storage: %w", err)
	}

	return nil
}

func (r *TeamRepository) GetByName(ctx context.Context, name string) (team string, err error) {
	const getByNameMethod = "GetByName"

	log := r.logger.WithGroup(getByNameMethod)

	operationLog := logger.StartOperation(
		log,
		"getting team from store",
		slog.String("team_name", name),
	)

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.String("team_name", name),
			slog.String("result_team_name", team),
		)
	}()

	team, err = r.queries.GetTeam(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.Join(apperrors.ErrTeamNotFound, err)
		}

		return "", fmt.Errorf("failed get team from postgres storage: %w", err)
	}

	return team, nil
}

func (r *TeamRepository) List(ctx context.Context) (teams []string, err error) {
	const listMethod = "List"

	log := r.logger.WithGroup(listMethod)

	operationLog := logger.StartOperation(log, "listing teams from store")

	defer func() {
		operationLog.FinishOperation(
			&err,
			slog.Int("teams_count", len(teams)),
		)
	}()

	teams, err = r.queries.ListTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed list teams from postgres storage: %w", err)
	}

	return teams, nil
}
