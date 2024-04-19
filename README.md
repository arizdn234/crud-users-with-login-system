# CRUD Users with Login System (Local Auth)

This project is a comprehensive **user management system** that includes features such as **login**, **registration**, **logout**, and **CRUD operations for user data**. It implements **authorization** for certain routes **using JSON Web Tokens (JWT)**, providing secure access to protected endpoints.

When a user **successfully logs in**, a **JWT is generated** and **set in the cookies**, with an **expiration time** of **1 day**. This ensures that authenticated users can access protected routes without having to repeatedly log in.

**User data management involves hashing passwords** for security purposes, ensuring that **sensitive information** is **securely stored** in the database. The project uses **SQLite as the database engine** and leverages the ORM framework **GORM for database operations**, simplifying database interactions and improving code readability.

To ensure code quality and reliability, the project **includes comprehensive testing** using **the Ginkgo** and **Gomega testing frameworks**. This helps identify and fix bugs early in the development process, ensuring a robust and stable user management system.

Overall, this project provides a scalable and secure solution for managing user data, with features such as authentication, authorization, and data encryption, making it suitable for a wide range of web applications and services.

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
