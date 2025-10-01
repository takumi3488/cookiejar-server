-- name: UpsertCookies :exec
INSERT INTO cookies (host, cookies) VALUES ($1, $2)
ON CONFLICT (host) DO UPDATE SET cookies = $2;

-- name: ListCookies :many
SELECT * FROM cookies;

-- name: DeleteAllCookies :exec
DELETE FROM cookies;