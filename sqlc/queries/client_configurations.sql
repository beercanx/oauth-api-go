-- name: GetClientConfiguration :one
SELECT *
FROM client_configurations
WHERE client_id = ?
LIMIT 1;
