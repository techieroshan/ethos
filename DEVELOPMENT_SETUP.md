# Ethos Development Setup Guide

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25.4+
- Node.js 16+
- Git

### Start Development Environment

```bash
cd /Users/roshanshah/Projects/ethos

# Start all services with dev profile (includes Mailpit)
docker-compose --profile dev up -d

# View logs
docker-compose logs -f ethos-api
```

### Access Services

| Service | URL | Purpose |
|---------|-----|---------|
| **Frontend** | http://localhost:5173 | React UI |
| **Backend API** | http://localhost:8000 | REST API (Gin) |
| **Mailpit** | http://localhost:8025 | Email testing inbox |
| **PostgreSQL** | localhost:5432 | Database |
| **Redis** | localhost:6379 | Cache |

---

## Environment Configuration

### Local Development (`.env.local`)
- **Email Service**: Mailpit (SMTP to localhost:1025)
- **Database**: Local PostgreSQL (no SSL)
- **Debug Mode**: Enabled
- **OTEL**: Disabled

### Production (`.env.prod`)
- **Email Service**: Emailit API (production email)
- **Database**: Remote PostgreSQL (SSL required)
- **Debug Mode**: Disabled
- **OTEL**: Enabled

### Switching Environments
```bash
# For local development (default)
cp .env.local .env

# For production
cp .env.prod .env
# Then update CHANGE_ME_ values with real credentials
```

---

## Email Testing with Mailpit

### Sending Test Emails
1. Create an account or sign up at http://localhost:5173
2. Click "Resend verification code"
3. Check Mailpit inbox at http://localhost:8025
4. Verify email appears with code/link

### Mailpit Features
- **Web UI**: http://localhost:8025 - view all emails
- **SMTP**: localhost:1025 - backend connects here
- **API**: http://localhost:8025/api/v1 - programmatic access

---

## Backend Configuration

All configuration is loaded from `.env` file at startup:

| Variable | Dev Value | Purpose |
|----------|-----------|---------|
| `MAILPIT_ENABLED` | true | Enable Mailpit for dev |
| `MAILPIT_SMTP_HOST` | localhost | Mailpit SMTP server |
| `MAILPIT_SMTP_PORT` | 1025 | Mailpit SMTP port |
| `DATABASE_URL` | postgres://... | Database connection |
| `REDIS_URL` | redis://localhost:6379 | Cache connection |
| `JWT_ACCESS_SECRET` | dev-key | JWT signing key |
| `OTEL_ENABLED` | false | Observability disabled in dev |

---

## Project Structure

```
ethos/
├── cmd/api/
│   └── main.go              # Backend entry point
├── internal/
│   ├── auth/                # Authentication module
│   ├── feedback/            # Feedback features
│   ├── organization/        # Organization management
│   ├── config/config.go     # Configuration loader
│   └── database/            # Migrations
├── api/
│   └── routes.go            # Route registration
├── pkg/
│   ├── email/              # Email service abstractions
│   │   ├── mailpit/        # Mailpit SMTP client
│   │   ├── emailit/        # Emailit API client
│   │   └── checker/        # Email validation API
│   └── ...
├── ethos-ui/                # React frontend (Vite)
├── docker-compose.yml       # Dev environment (with profiles)
├── docker-compose.prod.yml  # Production environment
├── .env.local              # Local dev configuration
├── .env.prod               # Production configuration
└── Dockerfile              # Backend multi-stage build
```

---

## Implemented Endpoints (28 Total)

### Phase 1 Week 1 - Auth (3 endpoints)
- POST `/api/v1/auth/verify-email`
- POST `/api/v1/auth/change-password`
- POST `/api/v1/auth/setup-2fa`

### Phase 1 Week 2 - Organization (10 endpoints)
- CRUD operations for organizations
- Member management
- Settings management

### Phase 1 Week 3 - Moderation (5 endpoints)
- Appeal management
- Moderation actions
- History tracking

### Phase 2 Week 1 - Feedback (5 endpoints)
- Update/Delete feedback
- Update/Delete comments
- Analytics

### Phase 2 Week 2 - Feedback (5 endpoints)
- Search feedback
- Get trending feedback
- Pin/Unpin feedback
- Get feedback stats

---

## Development Commands

```bash
# Build backend
go build -o bin/ethos ./cmd/api

# Run tests
go test -v ./...

# Lint code
go vet ./...

# Build frontend
cd ethos-ui && npm run build

# Start dev environment
docker-compose --profile dev up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Database migrations (if configured)
make db-migrate

# Clean up volumes
docker-compose down -v
```

---

## Troubleshooting

### Mailpit not receiving emails
1. Check if `MAILPIT_ENABLED=true` in `.env`
2. Verify backend is using Mailpit: `docker logs ethos-api | grep Mailpit`
3. Check SMTP connection: `docker logs ethos-mailpit`

### Backend won't start
1. Check database is healthy: `docker ps | grep postgres`
2. View errors: `docker logs ethos-api`
3. Rebuild: `docker-compose --profile dev down && docker-compose --profile dev up -d`

### Frontend won't load
1. Check Nginx/frontend status: `docker logs ethos-frontend`
2. Clear browser cache
3. Check API is responding: `curl http://localhost:8000/api/v1/health`

### Database connection issues
1. Verify PostgreSQL is healthy: `docker logs ethos-postgres`
2. Check DATABASE_URL in `.env`
3. Test connection: `psql -U ethos_user -d ethos -h localhost`

---

## Next Steps

1. ✅ Start dev environment with Mailpit
2. ✅ Test email verification flow
3. ✅ Verify all 28 endpoints
4. → Phase 3: Notifications & Subscriptions
5. → Phase 4: Dashboard & Analytics

---

## Support

For issues or questions:
1. Check logs: `docker-compose logs -f`
2. Review configuration: `cat .env`
3. Test connectivity: `curl -v http://localhost:8000/api/v1/health`
