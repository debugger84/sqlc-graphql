-- name: GetPostComments :many
-- gql: Post.comments
SELECT * FROM "comment".comment
WHERE post_id = $1 LIMIT 1;

-- name: LeaveComment :one
-- gql: Mutation
INSERT INTO "comment".comment (id, comment, author_id, post_id, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: DeleteComment :exec
-- gql: Mutation
DELETE FROM "comment".comment
WHERE id = $1;
