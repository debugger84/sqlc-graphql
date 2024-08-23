# Base example of using the plugin

This example demonstrates how to use the plugin to generate GraphQL schema from the SQL schema.
It generates enum, the main structs, and the queries.

## How to use
1. Install sqlc (https://docs.sqlc.dev/en/latest/overview/install.html)
2. Add new queries to the storage/query.sql file
3. Mark all new queries with the `-- gql: <extended_type>.<query_name>` comment if you want to have its representation in GraphQL schema. For example:
```sql
-- gql: Query.getAllAuthors
SELECT * FROM authors;
```
Also, you can add only the `-- gql: <extended_type>` comment if you want to have the same name in Golang and GraphQL.
```sql 
-- gql: Query
SELECT * FROM authors;
```

3. Run the following commands to generate golang and graphql code:
```bash
cd examples/simple 
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
func (r *queryResolver) GetAllAuthors(ctx context.Context) ([]storage.Author, error) {
    panic(fmt.Errorf("not implemented"))
}
```

7. Change the panic to the real implementation of the query. For example:
```go
func (r *queryResolver) GetAllAuthors(ctx context.Context) ([]storage.Author, error) {
    return r.Queries.GetAllAuthors(ctx)
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
    listAuthors {
        id
        name
        bio
    }
}
```
and you will see the list of authors like this:
```json
{
  "data": {
    "listAuthors": [
      {
        "id": "1",
        "name": "John Doe",
        "bio": "Unknown guy"
      }
    ]
  }
}
```