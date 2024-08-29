-- name: GetPostComments :many
-- gql: Post.comments
SELECT * FROM "comment".comment
WHERE post_id = $1 LIMIT @count OFFSET @after;

-- name: LeaveComment :one
-- gql: Mutation
INSERT INTO "comment".comment (id, comment, author_id, post_id, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: DeleteComment :exec
-- gql: Mutation.deleteComment
DELETE FROM "comment".comment
WHERE id = $1 and author_id = $2;
