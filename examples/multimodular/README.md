# Multi-modular project example

This example demonstrates how to use sqlc with multiple plugins in a multi-modular project.
Multiple modules means that you have a project with several packages that are located in different directories splitting the project to bounded contexts. 
This kind of project is useful when you have a big project with a lot of entities and queries, and you want to organize them to reduce code complexity. Sometimes this type of project is called "Modular Monolith". 
Usually in modular monoliths, each module has its own database schema and its own set of queries.

This example is a very simple blog with a few lines of code. 
Its purpose is only show how to split your big application to several modules.

In this example, we have two modules: `post` and `comment`. 
Each module has its own database schema with the appropriate schema names `post` and `comment`.
Each module has its own queries and models.

If you want to use the default naming strategy of SQLc when it adds a schema name to an entity, if the schema is not public, then use the default
https://github.com/sqlc-dev/sqlc-gen-go generation plugin. 
In this case in each module you will have such entities like PostPost CommentComment.

But if you want to have inside a module all entities without such namespacing, please, use
https://github.com/debugger84/sqlc-gen-go generation plugin instead. It is a fork of original one with the new configuration option "default_schema".
Type this option for each module with appropriate values and enjoy with readable names like Post and Comment.

## Run the example
1. Install sqlc (https://docs.sqlc.dev/en/latest/overview/install.html)
2. Checkout repository and build a plugin to the bin folder with the following command:
```shell
make all
```
3. Go to the `examples/multimodular` directory
4. Run `make generate` to generate the GraphQL schemas and SQL queries.
5. Edit `.env` file and set your database connection string
6. Run `make migrate` to create the database schema
7. Run `make run` to start the server
8. Open the browser and go to the http://localhost:8082/. You will see the GraphQL playground.
9. Try to add a post to the database:
```graphql
mutation {
    makeDraft(request: {
        title: "Post 1",
        content: "Text 1"
    }
    ) {
        id
        title
        content
        authorId
    }
}
```
Run it with the header `Authorization Bearer <uuid4>` to make a draft with the author id. UUID4 is any random UUID4 string.
In playground it will look like this:
```json
{"Authorization": "Bearer da25a384-20d1-46ee-883b-16b86018a4c2"}
```
in the Headers section.
10. Publish the post sending the id obtained from the previous step:
```graphql
mutation {
    publish(
        id: "1ef65ffd-e660-6fd2-a3fe-4124aa42c3eb"
    ) {
        id
        title
        content
        status
    }
}
```
11. Add a comment to the post:
```graphql
mutation {
    leaveComment(request: {
        postId: "1ef65ffd-e660-6fd2-a3fe-4124aa42c3eb",
        comment: "Comment 1"
    }) {
        id
        postId
        comment
    }
}
```
12. Get the post with comments:
```graphql
{
    post(id: "1ef65ffd-e660-6fd2-a3fe-4124aa42c3eb") {
        id
        title
        content
        status
        comments(request:{count:10, after:0}) {
            id
            postId
            comment
        }
    }
}
```