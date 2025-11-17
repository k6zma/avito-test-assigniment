package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/team"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/user"
)

type TeamHandler struct {
	teamService *team.TeamService
	userService *user.UserService
}

func NewTeamHandler(
	teamService *team.TeamService,
	userService *user.UserService,
) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		userService: userService,
	}
}

func (h *TeamHandler) AddTeam(c *fiber.Ctx) error {
	var dto generated.Team

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(apiError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	teamDomain, err := mappers.ToDomainTeam(dto)
	if err != nil {
		return MapError(c, err)
	}

	if err = h.teamService.CreateTeam(c.Context(), teamDomain.Name); err != nil {
		return MapError(c, err)
	}

	for _, m := range teamDomain.Members {
		_, err = h.userService.UpsertUser(
			c.Context(),
			m.ID,
			m.Name,
			teamDomain.Name,
			m.IsActive,
		)
		if err != nil {
			return MapError(c, err)
		}
	}

	updatedTeam, err := h.teamService.GetTeam(c.Context(), teamDomain.Name)
	if err != nil {
		return MapError(c, err)
	}

	return c.Status(201).JSON(struct {
		Team generated.Team `json:"team"`
	}{
		Team: mappers.ToDTOTeam(updatedTeam),
	})
}

func (h *TeamHandler) GetTeam(
	c *fiber.Ctx,
	params generated.GetTeamParams,
) error {
	teamDomain, err := h.teamService.GetTeam(c.Context(), params.TeamName)
	if err != nil {
		return MapError(c, err)
	}

	teamDTO := mappers.ToDTOTeam(teamDomain)

	return c.Status(200).JSON(teamDTO)
}
