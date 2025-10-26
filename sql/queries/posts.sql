-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
);

-- name: GetPostsByUser :many
SELECT p.*
FROM posts p
JOIN feed_follows ff ON p.feed_id = ff.feed_id
JOIN users u ON ff.user_id = u.id
WHERE u.name = $1
LIMIT $2;