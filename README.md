# sqlc-graphql
This is a plugin for the sqlc compiler (https://sqlc.dev/). 
It adds the generation of GraphQL schema from the database schema.
This plugin was created using the codebase of sqlc-gen-go plugin (https://github.com/sqlc-dev/sqlc-gen-go).

It has the capability to create a GraphQL schema based on the database schema, as well as to derive GraphQL queries from SQL queries. The produced code is compatible with the gqlgen library (https://gqlgen.com/).

## Features
- Generates GraphQL schema from the database schema
- Generates GraphQL enums
- Generates comments for the GraphQL queries
- Generates queries for the GraphQL schema using the SQL queries as a base.

## How to use
1. Install sqlc (https://docs.sqlc.dev/en/latest/overview/install.html)
2. Make a project with sqlc (https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
3. Change the sqlc.yaml file to use plugins instead of direct generation
```yaml
version: '2'
plugins:
  - name: graphql
    wasm:
      url: https://github.com/debugger84/sqlc-graphql/releases/download/v0.1.2/sqlc-graphql.wasm
      sha256: fe1bbdf7679a24c18cdbef48d2519c9ef517db6316d264f4857a21ef8d3b4e9f

  - name: golang
    wasm:
      url: https://downloads.sqlc.dev/plugin/sqlc-gen-go_1.3.0.wasm
      sha256: e8206081686f95b461daf91a307e108a761526c6768d6f3eca9781b0726b7ec8
sql:
  - engine: "postgresql"
    schema: "schema.sql"
    queries: "query.sql"
    
    codegen:
      - plugin: graphql
        out: "./graphql"
        options:
          ## the full package name pointing to the code generated by golang plugin
          package: "tutorial/tutorial"
          ## generate GraphQL enums
          emit_all_enum_values: true
          ## create several default types and directives to work in conjunction with the gqlgen library https://gqlgen.com/
          gen_common_parts: true
          ## override a column type with a custom GraphQL type
          ## the type should be described manually in the extended.graphql file 
          overrides:
            - column: "test.img"
              gql_type: "Image"
              nullable: true
          ## exclude columns from the generated schema
          ## Test - is the generated Graphql object 
          ## and CreatedAt is the column name to be excluded    
          exclude:
            - "Test.CreatedAt"
      ## options for the default golang generation plugin https://github.com/sqlc-dev/sqlc-gen-go
      - plugin: golang
        out: "./"
        options:
          package: "tutorial"
          sql_package: "pgx/v4"
          emit_json_tags: true
          emit_all_enum_values: true
          json_tags_case_style: "camel"
          out: "./"

          overrides:
            - column: "test.img"
              go_type: "tutorial/tutorial.NullImage"
              nullable: true

```

4. Run the sqlc command
```bash
sqlc -f ./tutorial/sqlc.yaml generate
```

It will generate the GraphQL schema and the queries in the graphql folder.
Also generated types will be linked to the generated by sqlc golang code.
After that you can generate the resolvers using gqlgen library (https://gqlgen.com/), and link resolvers to SQL queries.

Examples of output:
```graphql
# The generated schema
enum AuthorStatus  @goModel(model: "simple/storage.AuthorStatus") {
    active
    inactive
    deleted
}

type Author @goModel(model: "simple/storage.Author") {
    id: Int!
    name: String!
    bio: String
    status: AuthorStatus!
}
```

```graphql
# The generated queries
extend type Mutation {
    createAuthor(request: CreateAuthorInput!): Author!
}

input CreateAuthorInput @goModel(model: "simple/storage.CreateAuthorParams") {
    name: String!
    bio: String
}
```



See the [examples](https://github.com/debugger84/sqlc-graphql/tree/main/examples) folder for more information.