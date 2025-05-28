# Users API

A RESTful API service for user management built with Go, PostgreSQL, and Docker. This service provides CRUD operations for user data with a clean architecture pattern.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete users
- **Clean Architecture**: Organized into handler, service, and repository layers
- **Database Migrations**: Automated database schema management with Goose
- **Input Validation**: Request validation using go-playground/validator
- **Containerized**: Docker and Docker Compose support
- **Graceful Shutdown**: Proper server shutdown handling
- **Unit Tests**: Comprehensive test coverage with mocks

## Tech Stack

- **Language**: Go 1.23
- **Database**: PostgreSQL 14
- **Router**: Gorilla Mux
- **Database Driver**: pgx/v5
- **Migration Tool**: Goose
- **Validation**: go-playground/validator
- **Testing**: Testify
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── handler/
│   │   └── user_handler.go    # HTTP handlers
│   ├── model/
│   │   └── user.go           # Data models and DTOs
│   ├── repository/
│   │   ├── user_repo.go      # Database operations
│   │   └── user_repository_test.go
│   └── service/
│       ├── user_service.go   # Business logic
│       └── user_service_test.go
├── migrations/
│   └── 20250408125326_create_sources.sql
├── docker-compose.yml
├── Dockerfile
├── Dockerfile.migrations
├── go.mod
└── go.sum
```

## API Endpoints

| Method | Endpoint     | Description    | Request Body |
|--------|-------------|----------------|--------------|
| POST   | `/users`    | Create user    | UserCreate   |
| GET    | `/users/{id}` | Get user by ID | -            |
| PUT    | `/users/{id}` | Update user    | UserUpdate   |
| DELETE | `/users/{id}` | Delete user    | -            |

### Request/Response Models

**User Model:**
```json
{
  "id": "uuid",
  "first_name": "string",
  "last_name": "string", 
  "email": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**UserCreate (POST /users):**
```json
{
  "first_name": "John",
  "last_name": "Doe", 
  "email": "john.doe@example.com"
}
```

**UserUpdate (PUT /users/{id}):**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@example.com"
}
```

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd users-api
   ```

2. **Start the services**
   ```bash
   docker-compose up --build
   ```

3. **API will be available at**: `http://localhost:8080`

### Manual Setup

1. **Prerequisites**
   - Go 1.23+
   - PostgreSQL 14+

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables**
   ```bash
   export DATABASE_URL="postgres://username:password@localhost:5432/dbname?sslmode=disable"
   export PORT="8080"
   ```

4. **Run migrations**
   ```bash
   goose -dir=./migrations postgres "$DATABASE_URL" up
   ```

5. **Start the server**
   ```bash
   go run cmd/main.go
   ```

## Configuration

The application uses environment variables for configuration:

| Variable     | Description           | Default |
|-------------|-----------------------|---------|
| DATABASE_URL | PostgreSQL connection string | Required |
| PORT        | Server port           | 8080    |

## Example Usage

**Create a user:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com"
  }'
```

**Get a user:**
```bash
curl http://localhost:8080/users/{user-id}
```

**Update a user:**
```bash
curl -X PUT http://localhost:8080/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "email": "jane.doe@example.com"
  }'
```

**Delete a user:**
```bash
curl -X DELETE http://localhost:8080/users/{user-id}
```

## Testing

**Run all tests:**
```bash
go test ./...
```

**Run tests with coverage:**
```bash
go test -cover ./...
```

**Run specific package tests:**
```bash
go test ./internal/service/
```

### Test Database Setup

For integration tests, set the `TEST_DATABASE_URL` environment variable:
```bash
export TEST_DATABASE_URL="postgres://username:password@localhost:5432/test_db?sslmode=disable"
```

## Database Migrations

The project uses [Goose](https://github.com/pressly/goose) for database migrations.

**Create a new migration:**
```bash
goose -dir=./migrations create migration_name sql
```

**Run migrations:**
```bash
goose -dir=./migrations postgres "$DATABASE_URL" up
```

**Rollback migrations:**
```bash
goose -dir=./migrations postgres "$DATABASE_URL" down
```

## Docker Services

The `docker-compose.yml` defines three services:

- **app**: The main Go application
- **db**: PostgreSQL database
- **migrations**: Migration runner service

**Useful Docker commands:**
```bash
# View logs
docker-compose logs app

# Access database
docker-compose exec db psql -U postgres -d users

# Rebuild and restart
docker-compose up --build

# Stop services
docker-compose down
```

## Architecture

The application follows a clean architecture pattern:

- **Handler Layer**: HTTP request/response handling
- **Service Layer**: Business logic and validation
- **Repository Layer**: Data access and database operations
- **Model Layer**: Data structures and DTOs

This separation ensures:
- Better testability with dependency injection
- Clear separation of concerns
- Easy to maintain and extend
- Database-agnostic business logic

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.