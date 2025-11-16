package errors

import "errors"

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")

	ErrUserNotFound = errors.New("user not found")
	ErrUserInactive = errors.New("user inactive")

	ErrPRNotFound       = errors.New("pull request not found")
	ErrPRExists         = errors.New("pull request already exists")
	ErrPRMerged         = errors.New("pull request already merged")
	ErrPRNotOpen        = errors.New("pull request not open")
	ErrReviewerNotFound = errors.New("reviewer not found")
	ErrNotAssigned      = errors.New("reviewer is not assigned")
	ErrNoCandidate      = errors.New("no candidate for reassignment")

	ErrValidation = errors.New("validation failed")
)
