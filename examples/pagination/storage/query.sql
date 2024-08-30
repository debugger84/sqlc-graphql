-- name: GetAuthor :one
-- gql: Query
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
-- gql: Query
-- paginated: offset
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
-- gql: Mutation
INSERT INTO authors (
    name, bio
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: DeleteAuthor :exec
-- gql: Mutation
DELETE FROM authors
WHERE id = $1;

-- name: UpdateAuthor :one
-- gql: Mutation
UPDATE authors
set name = $2,
    bio = $3
WHERE id = $1
RETURNING *;

-- name: DeactivateAuthor :one
-- gql: Mutation
UPDATE authors
set status = 'inactive'
WHERE id = $1
RETURNING *;