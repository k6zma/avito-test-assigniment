-- name: CreatePullRequest :exec
INSERT INTO pull_requests (
    pull_request_id, pull_request_name, author_id, status, created_at
)
VALUES ($1, $2, $3, 'OPEN', NOW());

-- name: GetPullRequest :one
SELECT *
FROM pull_requests
WHERE pull_request_id = $1;

-- name: MergePullRequest :one
UPDATE pull_requests
SET status = 'MERGED',
    merged_at = COALESCE(merged_at, NOW())
WHERE pull_request_id = $1
    RETURNING *;

-- name: ListPRsAssignedToUser :many
SELECT pr.*
FROM pull_requests pr
    JOIN pull_request_reviewers r
        ON pr.pull_request_id = r.pull_request_id
WHERE r.reviewer_id = $1
ORDER BY pr.created_at DESC;
