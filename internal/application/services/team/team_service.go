package team

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
	"github.com/k6zma/avito-test-assigniment/pkg/logger"
)

type TeamService struct {
	teams  repositories.TeamRepository
	users  repositories.UserRepository
	logger *slog.Logger
}

func NewTeamService(
	teamRepo repositories.TeamRepository,
	userRepo repositories.UserRepository,
	logger *slog.Logger,
) *TeamService {
	const teamEntity = "team"

	return &TeamService{
		teams:  teamRepo,
		users:  userRepo,
		logger: logger.WithGroup(teamEntity),
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string) (err error) {
	const createTeamMethod = "CreateTeam"

	log := s.logger.WithGroup(createTeamMethod)

	operationLog := logger.StartOperation(
		log,
		"creating team",
		slog.String("team_name", name),
	)
	defer operationLog.FinishOperation(&err, slog.String("team_name", name))

	if err = s.teams.Create(ctx, name); err != nil {
		return fmt.Errorf("failed create team: %w", err)
	}

	return nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (team *models.Team, err error) {
	const getTeamMethod = "GetTeam"

	log := s.logger.WithGroup(getTeamMethod)

	operationLog := logger.StartOperation(
		log,
		"getting team",
		slog.String("team_name", name),
	)
	defer operationLog.FinishOperation(&err, slog.String("team_name", name))

	if _, err := s.teams.GetByName(ctx, name); err != nil {
		return nil, fmt.Errorf("failed get team by name: %w", err)
	}

	userIDS, err := s.users.ListUsersInTeam(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed list users in team: %w", err)
	}

	members := make([]*models.TeamMember, 0, len(userIDS))

	for _, userID := range userIDS {
		user, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed get user %s for team %s: %w",
				userID,
				name,
				err,
			)
		}

		member, err := models.NewTeamMember(user.ID, user.Name, user.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed create team member domain model: %w", err)
		}

		members = append(members, member)
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("team %s has no active members", name)
	}

	team, err = models.NewTeam(name, members)
	if err != nil {
		return nil, fmt.Errorf("failed create team domain model: %w", err)
	}

	return team, nil
}

func (s *TeamService) ListTeams(ctx context.Context) (teams []string, err error) {
	const listTeamsMethod = "ListTeams"

	log := s.logger.WithGroup(listTeamsMethod)

	operationLog := logger.StartOperation(log, "listing teams")
	defer operationLog.FinishOperation(&err)

	teams, err = s.teams.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed list teams: %w", err)
	}

	return teams, nil
}
