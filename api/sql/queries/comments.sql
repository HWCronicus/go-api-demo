-- name: CreateComment :one
INSERT INTO comments (user_id, content)
VALUES ($1, $2)
RETURNING id, user_id, content, created_at, (SELECT email FROM users WHERE id = $1);

-- name: ListComments :many
SELECT c.*, u.email
FROM comments c
JOIN users u ON c.user_id = u.id
ORDER BY c.created_at DESC;
