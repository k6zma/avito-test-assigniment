-- name: CreateTeam :exec
INSERT INTO teams (team_name)
VALUES ($1)
    ON CONFLICT (team_name) DO NOTHING;

-- name: GetTeam :one
SELECT team_name
FROM teams
WHERE team_name = $1;

-- name: ListTeams :many
SELECT team_name
FROM teams
ORDER BY team_name;