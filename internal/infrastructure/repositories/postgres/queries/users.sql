-- name: UpsertUser :exec
INSERT INTO users (user_id, username, team_name, is_active)
VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id)
DO UPDATE SET
    username = EXCLUDED.username,
           team_name = EXCLUDED.team_name,
           is_active = EXCLUDED.is_active;

-- name: GetUser :one
SELECT *
FROM users
WHERE user_id = $1;

-- name: ListActiveUsersInTeam :many
SELECT user_id
FROM users
WHERE team_name = $1
  AND is_active = TRUE
ORDER BY user_id;

-- name: ListUsersInTeam :many
SELECT user_id
FROM users
WHERE team_name = $1
ORDER BY user_id;

-- name: SetUserActive :one
UPDATE users
SET is_active = $2
WHERE user_id = $1
    RETURNING *;
