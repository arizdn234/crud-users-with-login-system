# CRUD API for Hotel Room Availability
...

## project-starter
- project-starter base structure command:
```sh
mkdir -p cmd/server internal/handlers internal/models internal/repository internal/server db/migrations config
touch cmd/server/main.go \
      internal/handlers/user_handler.go \
      internal/models/user.go \
      internal/repository/user_repository.go \
      internal/server/server.go \
      go.mod

```
> This command creates the basic structure for project, including directories for server code, handlers, models, repository, server setup, database migrations, and configuration files.


## define models
> In the internal/models directory, define the `models` needed for application. This typically includes a User model to represent user data.

## define repository
> In the internal/repository directory, define the `repository layer` responsible for interacting with the database. Implement functions to perform `CRUD operations` on user data.

## define handlers
> In the internal/handlers directory, define the `HTTP request handlers` for application. These handlers will handle incoming HTTP requests, validate input, and call the appropriate repository functions.

## define server
> In the internal/server directory, define the `server setup and configuration`. This includes setting up the HTTP server, defining routes, middleware, and any other server-related configurations.

## define config
> In the config directory, define configuration files for application. This may include `environment-specific configuration files` (e.g., development.yaml, production.yaml) for managing different settings in different environments.

## define migration
> In the db/migrations directory, define `database migration files`. These files contain SQL statements to create, modify, or delete database schema objects.

## define test
> In the internal directory, alongside models, handlers, and repository code, define `tests` to ensure that your application functions correctly. Use testing frameworks like Go's built-in testing package or external packages like `Ginkgo` and `Gomega` for behavior-driven development.
