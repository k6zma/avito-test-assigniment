package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/configs"
)

func NewPGXPool(ctx context.Context, dbCfg *configs.PostgresDatabase) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed parse pgxpool config: %w", err)
	}

	poolCfg.MaxConns = int32(dbCfg.MaxConnections) //nolint:gosec
	poolCfg.MinConns = int32(dbCfg.MinConnections) //nolint:gosec
	poolCfg.ConnConfig.ConnectTimeout = dbCfg.ConnectionTimeout
	poolCfg.MaxConnLifetime = dbCfg.ConnectionTimeout
	poolCfg.MaxConnIdleTime = dbCfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = dbCfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed create pgxpool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, dbCfg.ConnectionTimeout)
	defer cancel()

	if err = pool.Ping(pingCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("failed ping database: %w", err)
	}

	return pool, nil
}
