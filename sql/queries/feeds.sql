-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT 
	feeds.*,
	users.name AS user_name
FROM 
	feeds
INNER JOIN
	users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT
	*
FROM
	feeds
WHERE
	$1 = feeds.url;

-- name: MarkFeedFetchedById :exec
UPDATE feeds
SET updated_at = $1, last_fetched_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT
	*
FROM
	feeds
ORDER BY 
	last_fetched_at ASC NULLS FIRST
LIMIT 1;

