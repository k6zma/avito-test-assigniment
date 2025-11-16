-- name: AddReviewer :exec
INSERT INTO pull_request_reviewers (id, pull_request_id, reviewer_id)
VALUES ($1, $2, $3)
    ON CONFLICT DO NOTHING;

-- name: ListReviewers :many
SELECT reviewer_id
FROM pull_request_reviewers
WHERE pull_request_id = $1;

-- name: RemoveReviewer :exec
DELETE FROM pull_request_reviewers
WHERE pull_request_id = $1 AND reviewer_id = $2;

-- name: CountReviewers :one
SELECT COUNT(*)
FROM pull_request_reviewers
WHERE pull_request_id = $1;

-- name: ActiveCandidatesForReassign :many
SELECT u.user_id
FROM users u
WHERE u.team_name = (
    SELECT u2.team_name FROM users u2 WHERE u2.user_id = $1
) AND u.is_active = TRUE AND u.user_id <> $1;
