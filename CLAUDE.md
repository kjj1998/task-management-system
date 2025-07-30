# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Development Setup
- Start development database: `make start-dev-db`
- Run migrations: `make migrate-up`
- Run application in development: `make run-dev`

### Testing
- Run all repository tests: `make repository-tests`
- Run specific test: `go test ./internal/repository/task -v`
- Run tests with coverage: `go test -cover ./internal/repository/...`

### Database Migrations
- Apply all migrations: `make migrate-up`
- Apply one migration: `make migrate-up-once`
- Rollback one migration: `make migrate-down`

### Code Quality
- Run linter: `make lint`
- Format code: `go fmt ./...`
- Run vet: `go vet ./...`

## Architecture Overview

This is a Go-based task management REST API with a clean architecture following the repository pattern:

### Core Components
- **Models** (`internal/models/`): Domain entities (User, Category, Task) with typed enums for TaskStatus and TaskPriority
- **Repository Layer** (`internal/repository/`): Data access interfaces and MySQL implementations for each entity
- **Services** (`internal/services/`): Business logic layer
- **Handlers** (`internal/handlers/`): HTTP request handlers
- **Server** (`internal/server/`): HTTP server setup with routing

### Database Design
- MySQL database with migrations in `migrations/` directory
- Three main tables: users, categories, tasks
- Uses golang-migrate for database migrations
- Test containers for integration testing

### Key Architectural Patterns
- **Repository Pattern**: Clean separation between data access and business logic
- **Dependency Injection**: Store aggregates all repositories, injected into services
- **Error Handling**: Centralized database error handling via `internal/errors/`
- **Environment Configuration**: Uses .env files for development, environment variables for production

### Application Flow
1. `cmd/api/main.go` - Entry point, loads environment and starts server
2. `internal/server/server.go` - Initializes database connection, creates store, services, and handlers
3. `internal/store/store.go` - Aggregates all repositories into a single store
4. Handlers → Services → Store → Repositories → Database

### Testing Strategy
- Uses testcontainers for integration tests with real MySQL database
- Repository tests include full CRUD operations
- Test utilities in `internal/repository/testutils/`
- Separate test database configuration

### Configuration
- Uses structured configuration management via `internal/config/config.go`
- Requires .env file for development with environment variables:
  - `ENV` - Environment (dev/prod)
  - `SERVER_HOST`, `SERVER_PORT` - Server configuration
  - `DB_USER`, `DB_PASS`, `DB_HOST`, `DB_PORT`, `DB` - Database configuration
  - `LOG_LEVEL` - Logging level
- Uses docker-compose.yml for local MySQL instance
- Application runs on configurable port (default 8080) with `/api` prefix
- Configuration includes validation and environment-specific defaults
