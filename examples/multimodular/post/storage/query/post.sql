-- name: MakeDraft :one
-- gql: Mutation
INSERT INTO post.post (id, title, content, author_id, status, created_at)
VALUES ($1, $2, $3, $4, 'draft', now())
RETURNING *;

-- name: Publish :one
-- gql: Mutation
UPDATE post.post SET status = 'published' WHERE id = $1 RETURNING *;

-- name: GetPosts :one
-- gql: Query.posts
SELECT * FROM post.post WHERE author_id = $1 ORDER BY created_at DESC;