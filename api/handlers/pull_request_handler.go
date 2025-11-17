package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/api/mappers"
	"github.com/k6zma/avito-test-assigniment/internal/application/services/pullrequest"
)

type PullRequestHandler struct {
	pullRequestService *pullrequest.PullRequestService
}

func NewPullRequestHandler(
	pullRequestService *pullrequest.PullRequestService,
) *PullRequestHandler {
	return &PullRequestHandler{pullRequestService: pullRequestService}
}

func (h *PullRequestHandler) CreatePullRequest(c *fiber.Ctx) error {
	var body generated.CreatePullRequestJSONRequestBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(apiError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	pullRequestID, err := uuid.Parse(body.PullRequestId)
	if err != nil {
		return c.Status(400).JSON(apiError{
			"VALIDATION_ERROR",
			"invalid pull_request_id",
		})
	}

	authorID, err := uuid.Parse(body.AuthorId)
	if err != nil {
		return c.Status(400).JSON(apiError{
			"VALIDATION_ERROR",
			"invalid author_id",
		})
	}

	pullRequest, err := h.pullRequestService.CreatePullRequest(
		c.Context(),
		pullRequestID,
		body.PullRequestName,
		authorID,
	)
	if err != nil {
		return MapError(c, err)
	}

	return c.Status(201).JSON(mappers.ToDTOPullRequest(pullRequest))
}

func (h *PullRequestHandler) MergePullRequest(c *fiber.Ctx) error {
	var body generated.MergePullRequestJSONRequestBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(apiError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	pullRequestID, err := uuid.Parse(body.PullRequestId)
	if err != nil {
		return c.Status(400).JSON(apiError{
			"VALIDATION_ERROR",
			"invalid pull_request_id",
		})
	}

	pullRequest, err := h.pullRequestService.MergePullRequest(c.Context(), pullRequestID)
	if err != nil {
		return MapError(c, err)
	}

	return c.Status(200).JSON(mappers.ToDTOPullRequest(pullRequest))
}

func (h *PullRequestHandler) ReassignReviewer(c *fiber.Ctx) error {
	var body generated.ReassignReviewerJSONRequestBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(apiError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid request body",
		})
	}

	pullRequestID, err := uuid.Parse(body.PullRequestId)
	if err != nil {
		return c.Status(400).JSON(apiError{
			"VALIDATION_ERROR",
			"invalid pull_request_id",
		})
	}

	oldUserID, err := uuid.Parse(body.OldUserId)
	if err != nil {
		return c.Status(400).JSON(apiError{
			"VALIDATION_ERROR",
			"invalid old_user_id",
		})
	}

	if err = h.pullRequestService.ReassignReviewer(
		c.Context(),
		pullRequestID,
		oldUserID,
	); err != nil {
		return MapError(c, err)
	}

	return c.SendStatus(204)
}
