package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/k6zma/avito-test-assigniment/api/servers"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/team"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/user"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/configs"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/flags"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/generated"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/repositories/postgres/store"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/storages/postgres"
	"github.com/k6zma/avito-test-assigniment/pkg/logger"
	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

func main() {
	app := fx.New(
		fx.NopLogger,

		fx.Provide(
			func() (*flags.PeerlyFlags, error) {
				if err := validators.InitValidators(); err != nil {
					return nil, err
				}

				return flags.GetHubFlags()
			},

			func(fls *flags.PeerlyFlags) (*configs.Config, error) {
				return configs.LoadConfig(fls.Environment)
			},

			func(cfg *configs.Config) (*slog.Logger, error) {
				if err := logger.InitLogger(cfg.LoggingCfg.Level, cfg.LoggingCfg.Format); err != nil {
					return nil, err
				}

				return slog.Default(), nil
			},

			func(cfg *configs.Config, log *slog.Logger) (*pgxpool.Pool, error) {
				pool, err := postgres.NewPGXPool(context.Background(), cfg.PostgresCfg)
				if err != nil {
					log.Error("failed to create postgres pool", slog.Any("error", err))

					return nil, err
				}

				return pool, nil
			},
		),

		fx.Provide(
			func(pool *pgxpool.Pool) *generated.Queries {
				return generated.New(pool)
			},
			func(q *generated.Queries, log *slog.Logger) repositories.TeamRepository {
				repoLogger := log.WithGroup("repositories")

				return store.NewTeamRepository(q, repoLogger)
			},
			func(q *generated.Queries, log *slog.Logger) repositories.UserRepository {
				repoLogger := log.WithGroup("repositories")

				return store.NewUserRepository(q, repoLogger)
			},
			func(q *generated.Queries, log *slog.Logger) repositories.PullRequestRepository {
				repoLogger := log.WithGroup("repositories")

				return store.NewPullRequestRepository(q, repoLogger)
			},
		),

		fx.Provide(
			func(log *slog.Logger, teamRepo repositories.TeamRepository, userRepo repositories.UserRepository) *team.TeamService {
				serviceLogger := log.WithGroup("services")

				return team.NewTeamService(teamRepo, userRepo, serviceLogger)
			},
			func(log *slog.Logger, userRepo repositories.UserRepository, teamRepo repositories.TeamRepository) *user.UserService {
				serviceLogger := log.WithGroup("services")

				return user.NewUserService(userRepo, teamRepo, serviceLogger)
			},
			func(
				log *slog.Logger,
				prRepo repositories.PullRequestRepository,
				userRepo repositories.UserRepository,
			) *pullrequest.PullRequestService {
				serviceLogger := log.WithGroup("services")

				return pullrequest.NewPullRequestService(prRepo, userRepo, nil, nil, serviceLogger)
			},
		),

		fx.Provide(
			func(
				log *slog.Logger,
				userService *user.UserService,
				teamService *team.TeamService,
				prService *pullrequest.PullRequestService,
			) *servers.Server {
				serverLogger := log.WithGroup("server")

				return servers.NewServer(serverLogger, nil, userService, teamService, prService)
			},
		),

		fx.Invoke(
			func(cfg *configs.Config, fls *flags.PeerlyFlags, log *slog.Logger) {
				configLogger := log.WithGroup("config")

				configLogger.Debug(
					"Peerly config",
					slog.String("environment", fls.Environment),
					slog.Group("server_config",
						slog.String("Host", cfg.ServerCfg.Host),
						slog.Int("Port", cfg.ServerCfg.Port),
						slog.String("ReadTimeout", cfg.ServerCfg.ReadTimeout.String()),
						slog.String("WriteTimeout", cfg.ServerCfg.WriteTimeout.String()),
						slog.String("ShutdownTimeout", cfg.ServerCfg.ShutdownTimeout.String()),
						slog.Bool("Prefork", cfg.ServerCfg.Prefork),
					),
					slog.Group(
						"postgres_database_config",
						slog.String("Host", cfg.PostgresCfg.Host),
						slog.Int("Port", cfg.PostgresCfg.Port),
						slog.String("Name", cfg.PostgresCfg.Name),
						slog.String("User", cfg.PostgresCfg.User),
						slog.Int("MaxConnections", cfg.PostgresCfg.MaxConnections),
						slog.Int("MinConnections", cfg.PostgresCfg.MinConnections),
						slog.String(
							"ConnectionTimeout",
							cfg.PostgresCfg.ConnectionTimeout.String(),
						),
						slog.String("MaxConnLifetime", cfg.PostgresCfg.MaxConnLifetime.String()),
						slog.String("MaxConnIdleTime", cfg.PostgresCfg.MaxConnIdleTime.String()),
						slog.String(
							"HealthCheckPeriod",
							cfg.PostgresCfg.HealthCheckPeriod.String(),
						),
					),
					slog.Group("logging_config",
						slog.String("Level", cfg.LoggingCfg.Level),
						slog.String("Format", cfg.LoggingCfg.Format),
					),
				)
			},
		),

		fx.Invoke(
			func(
				lc fx.Lifecycle,
				srv *servers.Server,
				cfg *configs.Config,
				pool *pgxpool.Pool,
				log *slog.Logger,
				sd fx.Shutdowner,
			) {
				serverLogger := log.WithGroup("server")

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						serverLogger.Info(
							"Starting HTTP REST API server",
							slog.String("host", cfg.ServerCfg.Host),
							slog.Int("port", cfg.ServerCfg.Port),
						)

						go func() {
							if err := srv.Run(cfg.ServerCfg); err != nil {
								serverLogger.Error(
									"HTTP REST API server exited with error",
									slog.Any("error", err),
								)

								if shutdownErr := sd.Shutdown(); shutdownErr != nil {
									serverLogger.Error(
										"Failed to shutdown fx app",
										slog.Any("error", shutdownErr),
									)
								}
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						shutdownCtx, cancel := context.WithTimeout(
							ctx,
							cfg.ServerCfg.ShutdownTimeout,
						)
						defer cancel()

						serverLogger.Info("Shutting down HTTP REST API server")

						if err := srv.Shutdown(shutdownCtx); err != nil {
							serverLogger.Error(
								"Failed to shutdown HTTP REST API server",
								slog.Any("error", err),
							)
						}

						serverLogger.Info("Closing database pgx pool")
						pool.Close()

						return nil
					},
				})
			},
		),
	)

	app.Run()
}
