# Ethos BFF Architecture Documentation

## Overview

Ethos is a **Backend for Frontend (BFF)** service built with Go, designed to aggregate data from multiple Core Backend services and optimize payloads for frontend consumption. The BFF handles authentication, authorization, and provides a unified API interface for client applications.

## Architecture Pattern: Backend for Frontend (BFF)

### What is BFF?

The BFF pattern involves creating a dedicated backend service that:
- **Aggregates** data from multiple backend services
- **Optimizes** payloads specifically for frontend needs
- **Handles** authentication and authorization
- **Reduces** frontend complexity by providing tailored endpoints
- **Improves** performance by reducing round trips

### Why BFF for Ethos?

1. **Frontend Optimization**: Tailored responses reduce payload size and complexity
2. **Service Aggregation**: Combines data from multiple Core Backend services
3. **Security Boundary**: Centralized authentication/authorization
4. **Protocol Flexibility**: Uses REST for simple operations, gRPC for performance-critical paths
5. **Independent Deployment**: Can be deployed and scaled independently

## System Architecture

```
┌─────────────┐
│   Frontend  │ (Web/Mobile Applications)
└──────┬──────┘
       │ REST (JSON)
       │ HTTPS
       │
┌──────▼──────────────────────────────────────────────┐
│              Ethos BFF Service                      │
│  ┌──────────────────────────────────────────────┐  │
│  │  HTTP Handlers (Gin Framework)               │  │
│  │  - Authentication Middleware                 │  │
│  │  - Request Validation                        │  │
│  │  - Response Transformation                   │  │
│  └──────────────┬───────────────────────────────┘  │
│                 │                                   │
│  ┌──────────────▼───────────────────────────────┐  │
│  │  Service Layer (Business Logic)              │  │
│  │  - Auth Service                              │  │
│  │  - Profile Service                           │  │
│  │  - Feedback Service                          │  │
│  │  - Notification Service                      │  │
│  │  - Dashboard Service                         │  │
│  │  - People Service                            │  │
│  └──────────────┬───────────────────────────────┘  │
│                 │                                   │
│  ┌──────────────▼───────────────────────────────┐  │
│  │  Repository Layer (Data Access)              │  │
│  │  - PostgreSQL Repository                     │  │
│  │  - gRPC Clients (for Core Backend)           │  │
│  │  - REST Clients (for Core Backend)           │  │
│  └──────────────┬───────────────────────────────┘  │
└─────────────────┼──────────────────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
┌───────▼──────┐   ┌───────▼──────┐
│  PostgreSQL  │   │ Core Backend │
│  (BFF DB)    │   │  Services    │
│              │   │              │
│  - Users     │   │  - REST API  │
│  - Tokens    │   │  - gRPC API  │
│  - Profiles  │   │              │
│  - Feedback  │   │              │
│  - etc.      │   │              │
└──────────────┘   └──────────────┘
```

## Technology Stack

### Core Technologies

- **Language**: Go 1.25+
- **Web Framework**: Gin (REST API)
- **Database**: PostgreSQL (via pgx/v5)
- **Authentication**: JWT (access + refresh tokens)
- **Observability**: OpenTelemetry (tracing, metrics, logs)
- **Testing**: testify, go.uber.org/mock

### External Services

- **Email Validation**: Checker API (https://checker-docs.gitbook.io/checker-docs)
- **Email Sending**: Emailit API (https://docs.emailit.com/)
- **Local Email Testing**: Mailpit (SMTP server for development)

### Protocol Support

- **REST**: Primary protocol for client-facing APIs
- **gRPC**: Used for BFF ↔ Core Backend communication (performance-critical operations)
- **Protocol Buffers**: For gRPC service definitions

## Component Architecture

### Domain-Driven Design

The BFF is organized into **domains**, each following Clean Architecture principles:

```
internal/
├── auth/          # Authentication domain
├── profile/        # User profile domain
├── feedback/       # Feedback domain
├── notifications/  # Notifications domain
├── dashboard/      # Dashboard aggregation domain
├── people/         # People search domain
├── community/      # Community domain
└── account/        # Account/security domain
```

Each domain follows this structure:
```
domain/
├── handler/        # HTTP handlers (Gin)
├── service/        # Business logic layer
├── repository/     # Data access layer
└── model/          # Domain models
```

### Clean Architecture Layers

1. **Handler Layer** (`internal/*/handler/`)
   - HTTP request/response handling
   - Input validation
   - Error response formatting
   - Uses Gin framework

2. **Service Layer** (`internal/*/service/`)
   - Business logic
   - Orchestration
   - Transaction management
   - Interface-driven design

3. **Repository Layer** (`internal/*/repository/`)
   - Data persistence (PostgreSQL)
   - External service clients (gRPC/REST)
   - Query optimization
   - OpenTelemetry tracing

4. **Model Layer** (`internal/*/model/`)
   - Domain entities
   - Value objects
   - Data transfer objects

## REST vs gRPC Strategy

### When to Use REST

**REST is used for:**
- Simple CRUD operations
- Low-frequency operations
- Operations where human readability matters
- Authentication flows
- Account management

**Examples:**
- `POST /auth/login`
- `PUT /profile/me`
- `DELETE /profile/me`
- `GET /account/security-events`

### When to Use gRPC

**gRPC is used for:**
- High-frequency read operations
- Complex queries and aggregations
- Real-time/streaming scenarios
- Performance-critical paths
- BFF ↔ Core Backend communication

**Examples:**
- `GET /feedback/feed` (paginated, frequently accessed)
- `GET /dashboard` (aggregates multiple data sources)
- `GET /notifications` (real-time, potentially streaming)
- `GET /people/search` (complex queries)

### Protocol Selection

Protocol selection is **configuration-driven** via environment variables:

```go
// Example configuration
type Config struct {
    FeedbackProtocol string // "rest" or "grpc"
    DashboardProtocol string // "rest" or "grpc"
    // ...
}
```

**Default Behavior:**
- Development: REST (easier debugging)
- Production: gRPC for performance-critical endpoints

## Email Services Architecture

### Email Validation (Checker)

**Purpose**: Prevent fake account creation by detecting temporary/disposable emails

**Integration:**
- Called during user registration
- Validates email before account creation
- Returns error if temporary email detected

**Configuration:**
```go
type CheckerConfig struct {
    APIKey     string
    BaseURL    string
    Timeout    time.Duration
    Retries    int
}
```

### Email Sending (Emailit)

**Purpose**: Send templated emails for user communication

**Templates:**
- Email verification
- Password reset
- Account deletion confirmation
- Security alerts

**Configuration:**
```go
type EmailitConfig struct {
    APIKey     string
    BaseURL    string
    Timeout    time.Duration
    Retries    int
}
```

### Local Email Testing (Mailpit)

**Purpose**: Local development email testing without external API calls

**Integration:**
- SMTP server running locally
- Feature flag switches between Emailit (prod) and Mailpit (local)
- Captures all emails for inspection

**Configuration:**
```go
type MailpitConfig struct {
    SMTPHost   string
    SMTPPort   int
    FromEmail  string
}
```

## Database Architecture

### PostgreSQL Schema

**Core Tables:**
- `users` - User accounts
- `refresh_tokens` - JWT refresh tokens
- `user_preferences` - User preferences
- `account_deletions` - Scheduled account deletions
- `feedback_items` - Feedback posts
- `feedback_comments` - Comments on feedback
- `feedback_reactions` - Reactions to feedback
- `notifications` - User notifications
- `notification_preferences` - Notification settings
- `security_events` - Security event log
- `data_exports` - Data export requests

### Migration Strategy

- All migrations in `internal/database/migrations/`
- Up migrations: `XXX_description.up.sql`
- Down migrations: `XXX_description.down.sql`
- Versioned sequentially

## Authentication & Authorization

### JWT Token Strategy

**Access Tokens:**
- Short-lived (15 minutes default)
- Contains user ID
- Used for API authentication
- Stored in `Authorization: Bearer <token>` header

**Refresh Tokens:**
- Long-lived (14 days default)
- Stored in database (hashed)
- Used to obtain new access tokens
- Rotated on use

### Authentication Flow

```
1. User registers → Account created (email unverified)
2. Verification email sent → User clicks link
3. Email verified → User can login
4. Login → Access + Refresh tokens issued
5. API calls → Access token in header
6. Token expires → Use refresh token to get new access token
```

### Authorization Middleware

All protected endpoints use `middleware.AuthMiddleware()`:
- Extracts token from header
- Validates token
- Injects user ID into context
- Returns 401 if invalid/missing

## Observability

### OpenTelemetry Integration

**Tracing:**
- All HTTP requests traced
- Database queries traced
- External API calls traced (Checker, Emailit, gRPC)
- Distributed tracing across services

**Metrics:**
- Request latency
- Error rates
- Throughput
- Database connection pool metrics

**Logging:**
- Structured JSON logs
- Trace ID correlation
- Log levels: DEBUG, INFO, WARN, ERROR

### Export Destinations

- **Jaeger**: Distributed tracing
- **Prometheus**: Metrics collection
- **OpenTelemetry Collector**: Unified export

## Error Handling

### Standardized Error Responses

All errors follow this format:
```json
{
    "error": "Human-readable error message",
    "code": "ERROR_CODE"
}
```

### Error Codes

- `AUTH_INVALID_CREDENTIALS` - Invalid login credentials
- `AUTH_TOKEN_EXPIRED` - Token has expired
- `AUTH_TOKEN_INVALID` - Invalid token format
- `AUTH_EMAIL_UNVERIFIED` - Email not verified
- `VALIDATION_FAILED` - Input validation failed
- `USER_NOT_FOUND` - User does not exist
- `PROFILE_NOT_FOUND` - Profile does not exist
- `SERVER_ERROR` - Internal server error

## Security Considerations

### Input Validation

- All inputs validated at handler layer
- SQL injection prevention (parameterized queries)
- XSS prevention (output encoding)
- CSRF protection (token-based)

### Password Security

- Passwords hashed with bcrypt (cost 10)
- Never stored in plain text
- Never logged

### Token Security

- Tokens signed with HMAC-SHA256
- Refresh tokens stored hashed in database
- Token rotation on refresh
- Short-lived access tokens

### Rate Limiting

- Consider implementing distributed rate limiting (Redis)
- Per-user rate limits
- Per-endpoint rate limits

## Deployment Architecture

### Containerization

- Docker container for BFF service
- PostgreSQL container for database
- Mailpit container for local email testing

### Environment Configuration

- Development: Local PostgreSQL, Mailpit, REST for all operations
- Staging: Managed PostgreSQL, Emailit, gRPC for performance-critical
- Production: Managed PostgreSQL, Emailit, gRPC for performance-critical

### Scaling Considerations

- **Horizontal Scaling**: Stateless BFF can scale horizontally
- **Database Connection Pooling**: Configured per instance
- **gRPC Connection Pooling**: Reuse connections to Core Backend
- **Load Balancing**: Round-robin or least-connections

## Development Workflow

### TDD Approach

All development follows **Test-Driven Development (TDD)**:

1. **RED**: Write failing tests
2. **GREEN**: Implement minimal code to pass tests
3. **REFACTOR**: Improve code while keeping tests green

### Testing Strategy

- **Unit Tests**: Handler, service, repository layers
- **Integration Tests**: End-to-end API tests
- **Mock External Services**: Checker, Emailit, gRPC clients

### Code Quality

- `go fmt` - Code formatting
- `goimports` - Import organization
- `golangci-lint` - Linting
- `go test -cover` - Test coverage

## API Endpoints Overview

### Authentication Domain
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/me` - Get current user
- `DELETE /api/v1/auth/setup-2fa` - Disable 2FA

### Profile Domain
- `GET /api/v1/profile/me` - Get own profile
- `PUT /api/v1/profile/me` - Update profile
- `PATCH /api/v1/profile/me/preferences` - Update preferences
- `DELETE /api/v1/profile/me` - Delete account
- `GET /api/v1/profile/user-profile` - Search profiles

### Feedback Domain
- `GET /api/v1/feedback/feed` - Get feedback feed
- `GET /api/v1/feedback/:id` - Get feedback detail
- `GET /api/v1/feedback/:id/comments` - Get comments
- `POST /api/v1/feedback` - Create feedback
- `POST /api/v1/feedback/:id/comments` - Add comment
- `POST /api/v1/feedback/:id/react` - Add reaction
- `DELETE /api/v1/feedback/:id/react` - Remove reaction

### Notifications Domain
- `GET /api/v1/notifications` - Get notifications
- `GET /api/v1/notifications/preferences` - Get preferences
- `PUT /api/v1/notifications/preferences` - Update preferences

### Dashboard Domain
- `GET /api/v1/dashboard` - Get dashboard snapshot

### People Domain
- `GET /api/v1/people/search` - Search people
- `GET /api/v1/people/recommendations` - Get recommendations

### Community Domain
- `GET /api/v1/community/rules` - Get community rules

### Account Domain
- `GET /api/v1/account/security-events` - Get security events
- `GET /api/v1/account/export-data/:id/status` - Get export status

## Future Enhancements

### Planned Features

1. **gRPC Server Implementation**: Core Backend services as gRPC servers
2. **Caching Layer**: Redis for frequently accessed data
3. **Message Queue**: For async operations (email sending, notifications)
4. **GraphQL Support**: Alternative to REST for flexible queries
5. **WebSocket Support**: Real-time notifications

### Performance Optimizations

1. **Response Caching**: Cache frequently accessed endpoints
2. **Database Query Optimization**: Index optimization, query analysis
3. **Connection Pooling**: Optimize gRPC connection reuse
4. **Payload Compression**: Gzip compression for large responses

## References

- [APISpecs.md](./APISpecs.md) - API specification
- [.cursorrules](./.cursorrules) - Development guidelines
- [Checker API Docs](https://checker-docs.gitbook.io/checker-docs)
- [Emailit API Docs](https://docs.emailit.com/)

