package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/handlers"
	"github.com/k6zma/avito-test-assigniment/api/middlewares"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/team"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/user"
	"github.com/k6zma/avito-test-assigniment/internal/infrastructure/configs"
)

type Server struct {
	app         *fiber.App
	logger      *slog.Logger
	api         generated.ServerInterface
	userService *user.UserService
	teamService *team.TeamService
	prService   *pullrequest.PullRequestService
}

func NewServer(
	log *slog.Logger,
	api generated.ServerInterface,
	userService *user.UserService,
	teamService *team.TeamService,
	prService *pullrequest.PullRequestService,
) *Server {
	return &Server{
		logger:      log,
		api:         api,
		userService: userService,
		teamService: teamService,
		prService:   prService,
	}
}

func (s *Server) Run(cfg *configs.Server) error {
	s.app = fiber.New(fiber.Config{
		AppName:               cfg.AppName,
		StrictRouting:         true,
		Concurrency:           cfg.Concurrency,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		IdleTimeout:           cfg.IdleTimeout,
		Prefork:               cfg.Prefork,
		DisableStartupMessage: true,
	})

	s.app.Use(recover.New())
	s.app.Use(cors.New())
	s.app.Use(middlewares.NewHTTPLogger(s.logger))

	s.app.Static("/api/docs", "./api/docs")

	s.app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/api/docs/openapi.yml",
	}))

	userHandler := handlers.NewUserHandler(s.userService, s.prService)
	teamHandler := handlers.NewTeamHandler(s.teamService, s.userService)
	prHandler := handlers.NewPullRequestHandler(s.prService)

	generated.RegisterHandlers(s.app, struct {
		*handlers.UserHandler
		*handlers.TeamHandler
		*handlers.PullRequestHandler
	}{
		UserHandler:        userHandler,
		TeamHandler:        teamHandler,
		PullRequestHandler: prHandler,
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.app == nil {
		return nil
	}

	return s.app.ShutdownWithContext(ctx)
}
