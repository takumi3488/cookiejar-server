-- name: UpsertCookies :exec
INSERT INTO cookies (host, cookies, updated_at) VALUES ($1, $2, $3)
ON CONFLICT (host) DO UPDATE SET cookies = $2, updated_at = $3;

-- name: ListCookies :many
SELECT * FROM cookies;

-- name: GetCookiesByHost :one
SELECT * FROM cookies WHERE host = $1;