package team

import (
	"context"
	"fmt"

	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
	"github.com/k6zma/avito-test-assigniment/internal/domain/repositories"
)

type TeamService struct {
	teams repositories.TeamRepository
	users repositories.UserRepository
}

func NewTeamService(
	teamRepo repositories.TeamRepository,
	userRepo repositories.UserRepository,
) *TeamService {
	return &TeamService{
		teams: teamRepo,
		users: userRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string) error {
	if err := s.teams.Create(ctx, name); err != nil {
		return fmt.Errorf("failed create team: %w", err)
	}

	return nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	if _, err := s.teams.GetByName(ctx, name); err != nil {
		return nil, fmt.Errorf("failed get team by name: %w", err)
	}

	ids, err := s.users.ListActiveUsersInTeam(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed list users in team: %w", err)
	}

	members := make([]*models.TeamMember, 0, len(ids))

	for _, id := range ids {
		user, err := s.users.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf(
				"failed get user %s for team %s: %w",
				id,
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

	team, err := models.NewTeam(name, members)
	if err != nil {
		return nil, fmt.Errorf("failed create team domain model: %w", err)
	}

	return team, nil
}

func (s *TeamService) ListTeams(ctx context.Context) ([]string, error) {
	teams, err := s.teams.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed list teams: %w", err)
	}

	return teams, nil
}
