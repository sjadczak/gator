-- name: CreateFeedFollow :one
WITH new_follow AS (
    INSERT INTO
        feed_follows (
            id,
            created_at,
            updated_at,
            user_id,
            feed_id
        )
    VALUES
        (
            $1,
            $2,
            $3,
            $4,
            $5
        )
    RETURNING
        *
)
SELECT
    nf.id,
    u.name AS username,
    f.name AS feedname
FROM
    new_follow nf
    INNER JOIN users u ON nf.user_id = u.id
    INNER JOIN feeds f on nf.feed_id = f.id;

-- name: GetUserFeedFollows :many
SELECT
    ff.*,
    u.name AS username,
    f.name AS feedname
FROM
    feed_follows ff
    INNER JOIN users u ON ff.user_id = u.id
    INNER JOIN feeds f ON ff.feed_id = f.id
WHERE
    ff.user_id = (
        SELECT users.id
        FROM users
        WHERE users.name = $1
    );