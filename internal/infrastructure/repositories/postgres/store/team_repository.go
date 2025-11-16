package store

import (
	"context"

	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
)

type TeamRepository struct {
	queries *generated.Queries
}

func NewTeamRepository(q *generated.Queries) repositories.TeamRepository {
	return &TeamRepository{queries: q}
}

func (r *TeamRepository) Create(ctx context.Context, name string) error {
	return r.queries.CreateTeam(ctx, name)
}

func (r *TeamRepository) GetByName(ctx context.Context, name string) (string, error) {
	return r.queries.GetTeam(ctx, name)
}

func (r *TeamRepository) List(ctx context.Context) ([]string, error) {
	return r.queries.ListTeams(ctx)
}
