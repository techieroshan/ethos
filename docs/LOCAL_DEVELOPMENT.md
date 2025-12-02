# Local Development Guide

## Prerequisites

- Go 1.25+
- PostgreSQL 14+
- Docker (for Mailpit)
- Protocol Buffer compiler (for gRPC development)

## Setup

### 1. Clone and Install Dependencies

```bash
git clone <repo-url>
cd ethos
go mod download
```

### 2. Database Setup

```bash
# Start PostgreSQL
docker run -d \
  --name ethos-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ethos \
  -p 5432:5432 \
  postgres:14

# Run migrations (when migration tool is set up)
# migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/ethos?sslmode=disable" up
```

### 3. Email Testing Setup

```bash
# Start Mailpit
docker run -d \
  --name mailpit \
  -p 1025:1025 \
  -p 8025:8025 \
  axllent/mailpit
```

View emails at: http://localhost:8025

### 4. Environment Configuration

Create `.env` file:

```bash
# Server
SERVER_PORT=8000

# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/ethos?sslmode=disable
DB_MAX_CONNECTIONS=25

# JWT
JWT_ACCESS_SECRET=your-access-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=336h

# OpenTelemetry
OTEL_SERVICE_NAME=ethos-api
JAEGER_URL=http://localhost:14268/api/traces
OTEL_ENABLED=false

# Email Validation (Checker) - Optional
CHECKER_API_KEY=
CHECKER_BASE_URL=https://api.checker.com

# Email Sending - Use Mailpit for local dev
MAILPIT_ENABLED=true
MAILPIT_SMTP_HOST=localhost
MAILPIT_SMTP_PORT=1025
MAILPIT_FROM_EMAIL=noreply@ethos.test

# Email Sending - Emailit (production)
EMAILIT_API_KEY=
EMAILIT_BASE_URL=https://api.emailit.com

# gRPC (optional for local dev)
GRPC_ENABLED=false
GRPC_FEEDBACK_ENDPOINT=localhost:50051
GRPC_DASHBOARD_ENDPOINT=localhost:50052
GRPC_NOTIFICATIONS_ENDPOINT=localhost:50053
GRPC_PEOPLE_ENDPOINT=localhost:50054
```

## Running the Application

```bash
# Run the server
go run cmd/api/main.go

# Or build and run
go build -o ethos-api cmd/api/main.go
./ethos-api
```

Server starts on `http://localhost:8000`

## Testing

### Run All Tests

```bash
go test ./... -v
```

### Run Tests with Coverage

```bash
go test ./... -cover
```

### Run Specific Domain Tests

```bash
go test ./internal/auth/... -v
go test ./internal/feedback/... -v
```

## Development Workflow

### TDD Approach

1. **RED**: Write failing tests
2. **GREEN**: Implement minimal code to pass
3. **REFACTOR**: Improve while keeping tests green

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check imports
goimports -w .
```

## Database Migrations

Migrations are in `internal/database/migrations/`:

- Up migrations: `XXX_description.up.sql`
- Down migrations: `XXX_description.down.sql`

To add a new migration:
1. Create `XXX_description.up.sql`
2. Create `XXX_description.down.sql`
3. Run migration tool

## Email Testing

With Mailpit enabled:
1. Register a new user
2. Check Mailpit UI at http://localhost:8025
3. Verify email content and templates

## gRPC Development

### Generate Proto Code

```bash
# Install protoc and plugins first
make proto
```

### Test gRPC Clients

```bash
# Requires gRPC server running
go test ./pkg/grpc/client/... -v
```

## Debugging

### OpenTelemetry Tracing

1. Start Jaeger: `docker run -d -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one`
2. Set `OTEL_ENABLED=true` in `.env`
3. View traces at http://localhost:16686

### Logging

Structured JSON logs are emitted (can be configured for local development)

## Common Issues

### Database Connection Failed

- Verify PostgreSQL is running
- Check `DATABASE_URL` in `.env`
- Ensure database exists

### Email Not Sending

- Check Mailpit is running (local dev)
- Verify `MAILPIT_ENABLED=true`
- Check Mailpit UI for errors

### gRPC Connection Failed

- Verify gRPC server is running
- Check endpoint configuration
- Review connection timeout

## Next Steps

- Set up migration tool (migrate, golang-migrate, etc.)
- Configure CI/CD pipeline
- Set up local Jaeger for tracing
- Add development documentation for frontend integration

