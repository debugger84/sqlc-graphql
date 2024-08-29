-- name: MakeDraft :one
-- gql: Mutation
INSERT INTO post.post (id, title, content, author_id, status, created_at)
VALUES ($1, $2, $3, $4, 'draft', now())
RETURNING *;

-- name: Publish :one
-- gql: Mutation
UPDATE post.post SET status = 'published' WHERE id = $1 RETURNING *;

-- name: GetLastPosts :many
-- gql: Query.lastPosts
SELECT * FROM post.post WHERE status= 'published' ORDER BY created_at DESC LIMIT @count OFFSET @after;

-- name: GetMyDrafts :many
-- gql: Query.myDrafts
SELECT * FROM post.post WHERE author_id = $1 AND status = 'draft' ORDER BY created_at DESC LIMIT @count OFFSET @after;


-- name: GetPost :one
-- gql: Query.post
SELECT * FROM post.post WHERE id = $1;