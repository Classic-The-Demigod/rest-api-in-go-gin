# rest-api-in-go

A RESTful API built with Go and the Gin framework, featuring JWT authentication, SQLite persistence, and database migrations.

## Tech Stack

- **Go** with [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- **SQLite** via [go-sqlite3](https://github.com/mattn/go-sqlite3) — database
- **golang-migrate** — database migrations
- **golang-jwt** — JWT authentication
- **bcrypt** — password hashing
- **swaggo/swag** — API documentation
- **air** — live reload for development

## Project Structure

```
.
├── cmd/
│   ├── api/
│   │   ├── main.go         # App entrypoint, config
│   │   ├── server.go       # HTTP server setup
│   │   ├── routes.go       # Route definitions
│   │   ├── auth.go         # Auth handlers (register, login)
│   │   ├── events.go       # Event handlers
│   │   ├── middleware.go   # JWT auth middleware
│   │   └── context.go      # Gin context helpers
│   └── migrate/
│       ├── main.go         # Migration runner
│       └── migrations/     # SQL migration files
├── internal/
│   ├── database/
│   │   ├── models.go       # Model registry
│   │   ├── users.go        # User model & queries
│   │   ├── events.go       # Event model & queries
│   │   └── attendees.go    # Attendee model & queries
│   └── env/
│       └── env.go          # Environment variable helpers
├── docs/                   # Auto-generated swagger docs
├── data.db                 # SQLite database
└── .air.toml               # Air live reload config
```

## Getting Started

### Prerequisites

- Go 1.21+
- GCC (required for go-sqlite3 CGO)
- [MSYS2](https://www.msys2.org) on Windows (for GCC via MinGW)

### Installation

```bash
git clone https://github.com/your-username/rest-api-in-go
cd rest-api-in-go
go mod download
```

### Environment Variables

Create a `.env` file in the project root:

```env
PORT=8080
JWT_SECRET=your-secret-key
```

### Running Migrations

```bash
# Windows (PowerShell)
$env:CGO_ENABLED=1
go run ./cmd/migrate/main.go up

# Run down
go run ./cmd/migrate/main.go down
```

### Running the Server

```bash
# Development (with live reload)
air

# Production
$env:CGO_ENABLED=1
go run ./cmd/api
```

The server starts on `http://localhost:8080`.

## API Reference

Interactive docs are available at `http://localhost:8080/swagger/index.html` when the server is running.

### Auth

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register a new user | No |
| POST | `/api/v1/auth/login` | Login and get JWT token | No |

### Events

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/events` | Get all events | No |
| GET | `/api/v1/events/:id` | Get event by ID | No |
| POST | `/api/v1/events` | Create an event | Yes |
| PUT | `/api/v1/events/:id` | Update an event | Yes |
| DELETE | `/api/v1/events/:id` | Delete an event | Yes |

### Attendees

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/events/:id/attendees` | Get attendees for an event | No |
| GET | `/api/v1/attendees/:id/events` | Get events for an attendee | No |
| POST | `/api/v1/events/:id/attendees/:userId` | Add attendee to event | Yes |
| DELETE | `/api/v1/events/:id/attendees/:userId` | Remove attendee from event | Yes |

### Authentication

Protected routes require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

Get a token by calling the login endpoint.

## Example Requests

**Register**
```json
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "password",
  "name": "John Doe"
}
```

**Login**
```json
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password"
}
```

**Create Event**
```json
POST /api/v1/events
Authorization: Bearer <token>

{
  "name": "Go Conference 2026",
  "description": "A conference for Go developers",
  "date": "2026-06-15T09:00:00Z",
  "location": "Lagos, Nigeria"
}
```

## Generating Swagger Docs

```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

## Notes

- SQLite uses `?` as the query placeholder, not `$1`
- CGO must be enabled to build (`CGO_ENABLED=1`)
- On Windows, GCC must be available in PATH (via MSYS2 MinGW)