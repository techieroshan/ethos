#!/bin/bash

# Ethos Platform Deployment Script
# Usage: ./deploy.sh [environment] [action]
# Example: ./deploy.sh production deploy

set -e

ENVIRONMENT=${1:-development}
ACTION=${2:-deploy}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi

    # Check if .env file exists
    if [ ! -f "$PROJECT_ROOT/.env" ]; then
        log_warning ".env file not found. Creating from template..."
        cp "$PROJECT_ROOT/env.example" "$PROJECT_ROOT/.env"
        log_warning "Please edit .env file with your configuration before deploying."
        exit 1
    fi

    log_success "Prerequisites check passed"
}

# Build Docker images
build_images() {
    log_info "Building Docker images..."

    cd "$PROJECT_ROOT"

    # Build backend image
    log_info "Building backend image..."
    docker build -t ethos-backend:latest .

    # Build frontend image
    log_info "Building frontend image..."
    cd ethos-ui
    docker build -t ethos-frontend:latest .
    cd "$PROJECT_ROOT"

    log_success "Docker images built successfully"
}

# Deploy to specified environment
deploy() {
    log_info "Deploying to $ENVIRONMENT environment..."

    cd "$PROJECT_ROOT"

    case $ENVIRONMENT in
        development)
            COMPOSE_FILE="docker-compose.yml"
            ;;
        staging)
            COMPOSE_FILE="docker-compose.staging.yml"
            ;;
        production)
            COMPOSE_FILE="docker-compose.prod.yml"
            ;;
        *)
            log_error "Unknown environment: $ENVIRONMENT"
            log_info "Supported environments: development, staging, production"
            exit 1
            ;;
    esac

    # Stop existing containers
    log_info "Stopping existing containers..."
    docker-compose -f "$COMPOSE_FILE" down || true

    # Start containers
    log_info "Starting containers..."
    docker-compose -f "$COMPOSE_FILE" up -d --build

    # Wait for services to be healthy
    log_info "Waiting for services to be healthy..."
    sleep 30

    # Check health
    check_health

    log_success "Deployment to $ENVIRONMENT completed successfully"
}

# Check service health
check_health() {
    log_info "Checking service health..."

    # Check backend health
    if curl -f http://localhost:8000/api/v1/health &> /dev/null; then
        log_success "Backend is healthy"
    else
        log_warning "Backend health check failed"
    fi

    # Check frontend health
    if curl -f http://localhost:3000/health &> /dev/null; then
        log_success "Frontend is healthy"
    else
        log_warning "Frontend health check failed"
    fi
}

# Run database migrations
run_migrations() {
    log_info "Running database migrations..."

    cd "$PROJECT_ROOT"

    # For development environment
    if [ "$ENVIRONMENT" = "development" ]; then
        docker-compose exec api go run github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
            -path /app/internal/database/migrations \
            -database "postgres://ethos_user:ethos_password@postgres:5432/ethos?sslmode=disable" \
            up
    else
        log_warning "Migrations for $ENVIRONMENT environment need to be run manually"
    fi

    log_success "Database migrations completed"
}

# Backup database
backup_database() {
    log_info "Creating database backup..."

    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    BACKUP_FILE="ethos_backup_$TIMESTAMP.sql"

    cd "$PROJECT_ROOT"

    if [ "$ENVIRONMENT" = "development" ]; then
        docker-compose exec -T postgres pg_dump -U ethos_user ethos > "backups/$BACKUP_FILE"
        log_success "Database backup created: backups/$BACKUP_FILE"
    else
        log_warning "Database backup for $ENVIRONMENT needs to be done manually"
    fi
}

# Rollback deployment
rollback() {
    log_info "Rolling back deployment..."

    cd "$PROJECT_ROOT"

    # Stop current deployment
    docker-compose -f docker-compose.yml down

    # Start previous version (if using image tags)
    log_info "Starting previous version..."
    # Add rollback logic here

    log_success "Rollback completed"
}

# Show usage
usage() {
    echo "Ethos Platform Deployment Script"
    echo ""
    echo "Usage: $0 [environment] [action]"
    echo ""
    echo "Environments:"
    echo "  development  - Local development environment"
    echo "  staging      - Staging environment"
    echo "  production   - Production environment"
    echo ""
    echo "Actions:"
    echo "  deploy       - Deploy to specified environment"
    echo "  build        - Build Docker images only"
    echo "  migrate      - Run database migrations"
    echo "  backup       - Create database backup"
    echo "  health       - Check service health"
    echo "  rollback     - Rollback to previous version"
    echo ""
    echo "Examples:"
    echo "  $0 development deploy    # Deploy to development"
    echo "  $0 production build      # Build images for production"
    echo "  $0 staging migrate       # Run migrations on staging"
}

# Main script logic
main() {
    case $ACTION in
        deploy)
            check_prerequisites
            build_images
            deploy
            ;;
        build)
            check_prerequisites
            build_images
            ;;
        migrate)
            run_migrations
            ;;
        backup)
            backup_database
            ;;
        health)
            check_health
            ;;
        rollback)
            rollback
            ;;
        *)
            usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
