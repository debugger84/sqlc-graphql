# Pagination queries generation example

This is an example based on the "simple" example. It shows how to generate cursor-based (https://graphql.org/learn/pagination/) pagination queries using the sqlc tool.

Main differences from the "simple" example:
1. The `storage/query.sql` file contains the `ListAuthors` query without the `limit`, `offset`, and `order by` parameters.
But it has the `-- paginated: cursor:name,id` comment. It means that the `ListAuthors` query will be paginated by the `cursor` parameter that consists of the fields `name` and `id`.
```sql
-- name: ListAuthors :many
-- gql: Query
-- paginated: cursor:name,id
SELECT * FROM authors
ORDER BY name;
```
For the cursor-based pagination we have such requirements:
- The query should have the `-- paginated: cursor:<fields>` comment.
- The query should NOT have the `ORDER BY` clause. Because the order is defined by the `cursor` parameter.
- The query should NOT have the `LIMIT` and `OFFSET` clauses. Because the `LIMIT` is set by generator and `OFFSET` is not used at all in this type of pagination.
- If the cursor field is not unique, the query should have the `id` field in the last cursor parameter. The `id` field is used to make the order unique.
- If you want to make the order in the opposite direction, you should use the `-` sign before cursor field. For example:
```sql
-- paginated: cursor:-name,id
```
- The `cursor` parameter should be a string. The generator will parse the cursor string and extract the values of the fields from it.
- All fields in the `cursor` parameter should be in the returning type. 
- The generator will add the `ORDER BY` clause to the query with the fields from the `cursor` parameter.

3. The `storage/sqlc.yaml` uses the https://github.com/debugger84/sqlc-gen-go plugin to generate the pagination queries instead of the original one.

4. After running the generation, the generated code contains all necessary Relay-compatible types.
```graphql
type Author @goModel(model: "cursor/storage.Author") {
    id: Int!
    name: String!
    bio: String
    status: AuthorStatus!
}

type AuthorConnection @goModel(model: "cursor/storage.AuthorConnection") {
    edges: [AuthorEdge!]!
    pageInfo: PageInfo!
}

type AuthorEdge @goModel(model: "cursor/storage.AuthorEdge") {
    node: Author!
    cursor: String!
}


extend type Query {
    listAuthors(request: ListAuthorsInput!): AuthorConnection!
}
```

5. The `graph/schema.graphql` file contains the `ListAuthorsInput` input type that is used to pass the pagination parameters to the query.
```graphql
input ListAuthorsInput @goModel(model: "cursor/storage.ListAuthorsParams") {
    first: Int! @goField(name: "limit")
    after: String! @goField(name: "cursor")
}
```

6. The `graph/common.graphq` file contains the `PageInfo` type that is linked to the external `PageInfo` go struct to use the same type in GraphQL and Golang generators.


## How to use
1. Install sqlc (https://docs.sqlc.dev/en/latest/overview/install.html)
2. Add new queries to the storage/query.sql file
3. Mark all new queries with the comment: 
```sql
-- gql: <extended_type>.<query_name>
-- paginated: cursor:<comma separated list of fields>
``` 
For example:
```sql
-- gql: Query.getAllAuthors
-- paginated: cursor:-name,id
SELECT * FROM authors;
```

3. Run the following commands to generate golang and graphql code:
```bash
cd examples/cursor 
sqlc generate -f storage/sqlc.yaml
```
4. The generated code will be in the storage and graphql folders
5. Run the following command to generate the GraphQL resolvers by GraphQl schema:
```bash
go run github.com/99designs/gqlgen generate
```
It will generate the resolvers in the graph/resolver folder.

6. Find the changes in the graph/resolver folder. 
The changes will consist of the new queries with panics like this.
```go
func (r *queryResolver) GetAllAuthors(ctx context.Context, request storage.GetAllAuthorsParams) (storage.AuthorPage, error) {
    panic(fmt.Errorf("not implemented"))
}
```

7. Change the panic to the real implementation of the query. For example:
```go
func (r *queryResolver) GetAllAuthors(ctx context.Context, request storage.GetAllAuthorsParams) (storage.AuthorPage, error) {
    return r.Queries.GetAllAuthors(ctx, request)
}
```

8. Apply migrations to the database:
```bash
docker run --rm -it --network=host -v "$(pwd)/storage:/db" ghcr.io/amacneil/dbmate -u "postgres://postgres:foobar@localhost:5432/test?sslmode=disable" up
```

9. Run the following command to start the server:
```bash
PORT=8081 PG_URI="postgres://postgres:foobar@127.0.0.1:5432/test?sslmode=disable" go run server.go
```

10. Open the browser and go to the http://localhost:8081/. You will see the GraphQL playground.
11. Try to add several authors to the database:
```graphql
mutation {
    createAuthor(request:{name:"John Doe", bio:"Unknown guy"}) {
        id
        name
    }
}
```
12. Try to get all authors:
```graphql
query {
    listAuthors(request:{first:10, after:""}) {
        edges{node{id, name}, cursor}
        pageInfo{hasNextPage, endCursor}
    }
}

```
and you will see the list of authors like this:
```json
{
  "data": {
    "listAuthors": {
      "edges": [
        {
          "node": {
            "id": 2,
            "name": "Test Testov"
          },
          "cursor": "eyJuYW1lIjoiVGVzdCBUZXN0b3YiLCJpZCI6Mn0="
        }
      ],
      "pageInfo": {
        "hasNextPage": false,
        "endCursor": "eyJuYW1lIjoiVGVzdCBUZXN0b3YiLCJpZCI6Mn0="
      }
    }
  }
}
```