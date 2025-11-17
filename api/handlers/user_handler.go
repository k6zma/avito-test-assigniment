package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/user"
)

type UserHandler struct {
	userService        *user.UserService
	pullRequestService *pullrequest.PullRequestService
}

func NewUserHandler(
	userService *user.UserService,
	pullRequestService *pullrequest.PullRequestService,
) *UserHandler {
	return &UserHandler{
		userService:        userService,
		pullRequestService: pullRequestService,
	}
}

func (h *UserHandler) SetUserActive(c *fiber.Ctx) error {
	var body generated.SetUserActiveJSONRequestBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(apiError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	userID, err := uuid.Parse(body.UserId)
	if err != nil {
		return c.Status(400).JSON(apiError{
			"VALIDATION_ERROR",
			"invalid user_id",
		})
	}

	userDomain, err := h.userService.SetUserActive(c.Context(), userID, body.IsActive)
	if err != nil {
		return MapError(c, err)
	}

	resp := struct {
		User generated.User `json:"user"`
	}{
		User: mappers.ToDTOUser(userDomain),
	}

	return c.Status(200).JSON(resp)
}

func (h *UserHandler) GetUserAssignedReviews(
	c *fiber.Ctx,
	params generated.GetUserAssignedReviewsParams,
) error {
	userID, err := uuid.Parse(params.UserId)
	if err != nil {
		return c.Status(400).JSON(apiError{
			"VALIDATION_ERROR",
			"invalid user_id",
		})
	}

	pullRequests, err := h.pullRequestService.ListPullRequestsAssignedToUser(c.Context(), userID)
	if err != nil {
		return MapError(c, err)
	}

	result := make([]generated.PullRequestShort, 0, len(pullRequests))

	for _, pullRequest := range pullRequests {
		result = append(result, mappers.ToDTOPRShort(pullRequest))
	}

	resp := struct {
		UserID       string                       `json:"user_id"`
		PullRequests []generated.PullRequestShort `json:"pull_requests"`
	}{
		UserID:       userID.String(),
		PullRequests: result,
	}

	return c.Status(200).JSON(resp)
}
