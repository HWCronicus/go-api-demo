-- name: CreateComment :one
INSERT INTO comments (user_id, content)
VALUES ($1, $2)
RETURNING id, user_id, content, created_at, (SELECT email FROM users WHERE id = $1);

-- name: ListComments :many
SELECT c.*, u.email
FROM comments c
JOIN users u ON c.user_id = u.id
ORDER BY c.created_at DESC;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE comments.id = $1 
AND comments.user_id = (SELECT users.id FROM users WHERE users.email = $2);