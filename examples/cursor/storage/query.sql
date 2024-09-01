-- name: ListAuthors :many
-- gql: Query
-- paginated: cursor:name,id
SELECT * FROM authors;

-- name: CreateAuthor :one
-- gql: Mutation
INSERT INTO authors (
    name, bio
) VALUES (
             $1, $2
         )
RETURNING *;
