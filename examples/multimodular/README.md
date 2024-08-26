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

## How to use
1. Install sqlc (https://docs.sqlc.dev/en/latest/overview/install.html)
2. 
