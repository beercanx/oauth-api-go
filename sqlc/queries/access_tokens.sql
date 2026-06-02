-- name: CreateAccessToken :exec
INSERT INTO access_tokens (id, username, client_id, scopes, issued_at, expires_at, not_before)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetAccessToken :one
SELECT *
FROM access_tokens
WHERE id = ?
LIMIT 1;

-- name: DeleteAccessToken :exec
DELETE
FROM access_tokens
WHERE id = ?;

-- name: DeleteExpiredAccessTokens :exec
DELETE
FROM access_tokens
WHERE datetime(expires_at) <= CURRENT_TIMESTAMP;
