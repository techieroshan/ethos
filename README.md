# Ethos Platform

[![CI/CD Pipeline](https://github.com/your-org/ethos/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/your-org/ethos/actions/workflows/ci-cd.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/ethos)](https://goreportcard.com/report/github.com/your-org/ethos)
[![codecov](https://codecov.io/gh/your-org/ethos/branch/main/graph/badge.svg)](https://codecov.io/gh/your-org/ethos)

A comprehensive multi-tenant feedback management platform built with Go and React. Features enterprise-grade security, tenant isolation, and scalable architecture.

## 🚀 Features

### Core Functionality
- **Multi-Tenant Architecture** - Complete tenant isolation and management
- **User Feedback System** - Create, manage, and track feedback across organizations
- **Advanced Search & Discovery** - Powerful filtering and recommendation engine
- **Real-time Notifications** - Email and in-app notification system
- **Escalation & Appeals** - Structured feedback escalation workflow

### Administration
- **Platform Admin** - Global system management and analytics
- **Organization Admin** - Tenant-specific administration and moderation
- **Content Moderation** - Automated and manual content review
- **User Management** - Comprehensive user lifecycle management
- **Audit & Compliance** - Full audit trails and regulatory reporting

### Technical Features
- **Clean Architecture** - Domain-driven design with proper layering
- **Enterprise Security** - JWT authentication, RBAC, tenant isolation
- **High Performance** - Optimized queries, caching, and async processing
- **Comprehensive Testing** - 35+ E2E tests across 5 browser environments
- **Production Ready** - Docker containerization, CI/CD, monitoring

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React UI      │    │   Gin API       │    │   PostgreSQL    │
│   (Frontend)    │◄──►│   (Backend)     │◄──►│   (Database)    │
│                 │    │                 │    │                 │
│ - SPA Routing   │    │ - REST API      │    │ - Multi-tenant  │
│ - Tenant Context│    │ - JWT Auth      │    │ - Audit Logs    │
│ - Real-time UI  │    │ - Rate Limiting │    │ - Full-text     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │     Redis       │    │   Email/SMTP    │
                       │   (Cache)       │    │   (External)    │
                       └─────────────────┘    └─────────────────┘
```

## 🛠️ Tech Stack

### Backend
- **Go 1.21** - High-performance backend language
- **Gin** - HTTP web framework
- **PostgreSQL** - Primary database
- **Redis** - Caching and session storage
- **JWT** - Authentication tokens
- **OpenTelemetry** - Observability and tracing

### Frontend
- **React 18** - Modern UI framework
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS** - Utility-first CSS framework
- **Playwright** - E2E testing
- **Vite** - Fast build tool

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Local development
- **GitHub Actions** - CI/CD pipeline
- **Nginx** - Reverse proxy and load balancing
- **PostgreSQL** - Database migrations

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)
- PostgreSQL 15+ (optional, Docker provided)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/ethos.git
   cd ethos
   ```

2. **Set up environment**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

3. **Start development environment**
   ```bash
   make setup
   make dev
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8000
   - Mailpit (email testing): http://localhost:8025

### Production Deployment

1. **Build and deploy with Docker Compose**
   ```bash
   docker-compose -f docker-compose.prod.yml up --build -d
   ```

2. **Or use the Makefile**
   ```bash
   make deploy-production
   ```

## 🧪 Testing

### Run All Tests
```bash
make test-all
```

### Run Specific Test Suites
```bash
# Backend tests
make test

# Frontend tests
make test-frontend

# E2E tests
make test-e2e
```

### Test Coverage
- **Backend**: Unit tests with race detection and coverage reporting
- **Frontend**: Jest testing with coverage reporting
- **E2E**: Playwright tests across multiple browsers
- **Integration**: API and database integration tests

## 📊 API Documentation

### Authentication
```bash
# Login
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@company.com",
  "password": "password123"
}
```

### Multi-Tenant Operations
```bash
# Switch tenant context
POST /api/v1/auth/tenants/{tenant_id}/switch
Authorization: Bearer {token}

# List user tenants
GET /api/v1/auth/tenants
Authorization: Bearer {token}
```

### Feedback Management
```bash
# Create feedback
POST /api/v1/feedback
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Great work on the project!",
  "rating": 5,
  "recipient_id": "user-123"
}
```

## 🔒 Security

### Authentication & Authorization
- JWT-based authentication with refresh tokens
- Role-based access control (RBAC)
- Multi-tenant data isolation
- Rate limiting and DDoS protection

### Data Protection
- Database-level tenant isolation
- Encrypted sensitive data
- Secure API endpoints
- Input validation and sanitization

### Compliance
- Comprehensive audit logging
- GDPR compliance features
- SOC 2 ready architecture
- Regular security updates

## 📈 Monitoring & Observability

### Health Checks
```bash
# Application health
GET /api/v1/health

# Readiness check
GET /api/v1/ready
```

### Metrics
- **Prometheus** metrics collection
- **OpenTelemetry** distributed tracing
- **Structured logging** with correlation IDs
- **Performance monitoring** and alerting

### Logging
- JSON-formatted logs
- Request/response logging
- Error tracking with Sentry
- Audit trail logging

## 🚀 Deployment

### Environment Variables
```bash
# Required
ENV=production
DATABASE_URL=postgres://user:pass@host:5432/db
JWT_SECRET=your-secret-key
API_PORT=8000

# Optional
REDIS_URL=redis://host:6379
SMTP_HOST=smtp.gmail.com
SENTRY_DSN=your-sentry-dsn
```

### Docker Deployment
```bash
# Build images
docker-compose build

# Run in production
docker-compose -f docker-compose.prod.yml up -d

# Scale services
docker-compose up -d --scale api=3
```

### Kubernetes Deployment
```bash
# Apply manifests
kubectl apply -f k8s/

# Check status
kubectl get pods
kubectl get services
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Workflow
- Follow TDD (Red-Green-Refactor) approach
- Write comprehensive tests for new features
- Ensure all tests pass before merging
- Follow Go and React best practices
- Update documentation for API changes

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built following Clean Architecture principles
- Inspired by enterprise feedback management platforms
- Thanks to the Go and React communities

## 📞 Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/your-org/ethos/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/ethos/discussions)

---

**Ethos** - Enterprise-grade feedback management for modern organizations.
