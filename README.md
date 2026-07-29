# Cooking Site Project

![Go](...)
![License](...)

REST API powering a cooking recipe platform. Supports user authentication, recipe management, comments, and user ratings.

---

## Features

- User authentication
- Recipe management
- Recipe comments
- Recipe ratings
- JWT Authentication
- PostgreSQL support
- Structured logging
- OpenAPI documentation

---

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Database Driver:** pgx
- **Authentication:** JWT
- **Environment loading:** godotenv
- **Configuration:** caarlos0/env
- **Logging:** zerolog

---

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL

### Clone the repository

```bash
git clone https://github.com/Tihomir-Tinkov/cooking-site-project.git

cd cooking-site-project
```

### Install dependencies

```bash
go mod tidy
```

### Configure environment

Copy and configure the example environment file:

```bash
cp .env.example .env
```

### Run the application

```bash
go run ./cmd/server/main.go
```

The server will start at:

```
http://localhost:8080
```

---

## API documentation

OpenAPI specification:

- `/docs/openapi.yaml`
- Interactive Swagger UI: https://Tihomir-Tinkov.github.io/cooking-site-project

---

## Authentication

Protected routes require a JWT Bearer token.

Example:

```
Authorization: Bearer <your_token>
```

---

## Project Structure

```text
cooking-site-project/
├── cmd/
│   └── server/
│       └── main.go
├── docs/
│   ├── assets/
│   └── swagger-ui/
├── internal/
│   ├── app/
│   │   ├── adapters/
│   │   ├── controllers/
│   │   │   ├── middleware/
│   │   │   └── responses/
│   │   ├── gateway/
│   │   ├── models/
│   │   ├── ports/
│   │   ├── repositories/
│   │   ├── routes/
│   │   ├── services/
│   │   ├── app.go
│   │   └── registry.go
│   ├── cache/
│   ├── config/
│   ├── conn/
│   │   └── db/
│   └── logger/
├── migrations/
├── .env.example
├── .gitignore
├── README.md
├── api.http
├── go.mod
└── go.sum
```

---

## Architecture

This project follows a layered architecture inspired by Clean Architecture.

- controllers → HTTP handlers
- middleware → request processing
- services → business logic
- repositories → PostgreSQL access
- ports → interfaces
- adapters → external implementations

---

## Configuration

Configuration is loaded from environment variables with the `APP_` prefix.
For local development, values can be provided through a `.env` file.

---

## Environment Variables

| Variable                       | Default            | Required | Description                        |
| ------------------------------ | ------------------ | -------- | ---------------------------------- |
| `APP_ENV`                      | —                  | No       | Application environment            |
| `APP_API_PORT`                 | —                  | Yes      | HTTP server port                   |
| `APP_API_CORS_ALLOW_ORIGINS`   | —                  | No       | Allowed CORS origins               |
| `APP_ROUTES_PREFIX`            | `api`              | No       | API route prefix                   |
| `APP_ROUTES_LOGGER_MIDDLEWARE` | `true`             | No       | Enable request logging middleware  |
| `APP_DB_HOST`                  | —                  | Yes      | PostgreSQL host                    |
| `APP_DB_PORT`                  | —                  | Yes      | PostgreSQL port                    |
| `APP_DB_USER`                  | —                  | Yes      | PostgreSQL user                    |
| `APP_DB_PASS`                  | —                  | Yes      | PostgreSQL password                |
| `APP_DB_NAME`                  | —                  | Yes      | PostgreSQL database name           |
| `APP_DB_SSL_MODE`              | —                  | No       | Enable SSL connection              |
| `APP_DB_MAX_IDLE_CONNS`        | `10`               | No       | Maximum idle connections           |
| `APP_DB_MAX_OPEN_CONNS`        | `25`               | No       | Maximum open connections           |
| `APP_DB_CONN_MAX_LIFETIME`     | `0s`               | No       | Connection lifetime                |
| `APP_JWT_SECRET`               | —                  | Yes      | JWT signing secret                 |
| `APP_JWT_ISSUER`               | `cooking-site-api` | No       | JWT issuer                         |
| `APP_JWT_EXPIRATION`           | `24h`              | No       | Token lifetime                     |
| `APP_STORE_PATH`               | `storage`          | No       | Storage directory                  |
| `APP_LOG_ENCODING`             | `json`             | No       | Log encoding (`json` or `console`) |
| `APP_LOG_LEVEL`                | `info`             | No       | Log level                          |
| `APP_ARGON2_MEMORY`            | `65536`            | No       | Argon2 memory cost                 |
| `APP_ARGON2_ITERATIONS`        | `3`                | No       | Argon2 iterations                  |
| `APP_ARGON2_PARALLELISM`       | `4`                | No       | Argon2 parallelism                 |
| `APP_ARGON2_SALT_LENGTH`       | `16`               | No       | Salt length in bytes               |
| `APP_ARGON2_KEY_LENGTH`        | `32`               | No       | Derived key length in bytes        |

---

## Database Migrations

Run migrations:

```bash
goose -dir migrations postgres <connection_string> up
```

---

## Development

Format code:

```bash
go fmt ./...
```

Vet:

```bash
go vet ./...
```

---

## Logging

The application outputs structured logs with request information, errors, and startup events.

---

## License

This project is licensed under ...
