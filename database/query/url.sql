-- name: CreateURL :one
INSERT INTO urls (
    orignal_url,
    short_code,
    is_custom,
    expired_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: IsShortCodeAvaliable :one
SELECT NOT EXISTS (
    SELECT 1 FROM urls
    WHERE short_code = $1
) AS is_available;

-- name: GetURLByShortCode :one
SELECT *
FROM urls
WHERE short_code = $1 AND expired_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: DeleteExpiredURLs :exec
DELETE FROM urls
WHERE expired_at <= CURRENT_TIMESTAMP;