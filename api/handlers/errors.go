package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	apperrors "github.com/k6zma/avito-test-assigniment/internal/domain/errors"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func MapError(c *fiber.Ctx, err error) error {
	msg := sanitizeMessage(err.Error())

	c.Locals("error_cause", err)

	switch {
	case errors.Is(err, apperrors.ErrTeamExists):
		return c.Status(400).JSON(apiError{"TEAM_EXISTS", msg})
	case errors.Is(err, apperrors.ErrTeamNotFound):
		return c.Status(404).JSON(apiError{"TEAM_NOT_FOUND", msg})

	case errors.Is(err, apperrors.ErrUserNotFound):
		return c.Status(404).JSON(apiError{"USER_NOT_FOUND", msg})
	case errors.Is(err, apperrors.ErrUserInactive):
		return c.Status(409).JSON(apiError{"USER_INACTIVE", msg})

	case errors.Is(err, apperrors.ErrPRExists):
		return c.Status(409).JSON(apiError{"PR_EXISTS", msg})
	case errors.Is(err, apperrors.ErrPRNotFound):
		return c.Status(404).JSON(apiError{"PR_NOT_FOUND", msg})
	case errors.Is(err, apperrors.ErrPRMerged):
		return c.Status(409).JSON(apiError{"PR_MERGED", msg})
	case errors.Is(err, apperrors.ErrPRNotOpen):
		return c.Status(409).JSON(apiError{"PR_NOT_OPEN", msg})

	case errors.Is(err, apperrors.ErrReviewerNotFound):
		return c.Status(404).JSON(apiError{"REVIEWER_NOT_FOUND", msg})
	case errors.Is(err, apperrors.ErrNotAssigned):
		return c.Status(409).JSON(apiError{"NOT_ASSIGNED", msg})
	case errors.Is(err, apperrors.ErrNoCandidate):
		return c.Status(409).JSON(apiError{"NO_CANDIDATE", msg})

	case errors.Is(err, apperrors.ErrValidation):
		return c.Status(400).JSON(apiError{"VALIDATION_ERROR", msg})

	default:
		return c.Status(500).JSON(apiError{
			Code:    "INTERNAL_ERROR",
			Message: sanitizeMessage(err.Error()),
		})
	}
}

func sanitizeMessage(s string) string {
	s = strings.ReplaceAll(s, "\\n", " : ")
	s = strings.ReplaceAll(s, "\n", " : ")
	s = strings.ReplaceAll(s, " : ", ": ")

	return s
}
