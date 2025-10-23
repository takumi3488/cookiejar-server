-- name: UpsertCookies :exec
INSERT INTO cookies (host, cookies) VALUES ($1, $2)
ON CONFLICT (host) DO UPDATE SET cookies = $2;

-- name: ListCookies :many
SELECT * FROM cookies;

-- name: GetCookiesByHost :one
SELECT * FROM cookies WHERE host = $1;