# Base example of using the plugin

This example demonstrates how to use the plugin to generate GraphQL schema from the SQL schema.
It generates enum, the main structs, and the queries.

## How to use
1. Install sqlc (https://docs.sqlc.dev/en/latest/overview/install.html)
2. Run the following commands:
```bash
cd examples/simple 
sqlc generate -f storage/sqlc.yaml
```