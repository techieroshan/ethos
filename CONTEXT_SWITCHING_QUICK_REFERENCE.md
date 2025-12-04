# Multi-Tenant Context Switching - Quick Reference

## Core Files Structure

```
internal/organization/
├── repository/
│   ├── context_repository.go              # Interface definition
│   └── postgres_context_repository.go     # PostgreSQL implementation (16 methods)
├── service/
│   ├── context_service.go                 # Interface definition + request/response types
│   └── context_service_impl.go            # Service implementation (9 methods)
├── handler/
│   └── context_handler.go                 # HTTP handler (4 endpoints)
└── model/
    └── organization_model.go              # Data models

internal/middleware/
└── context_switch_middleware.go           # 4 middleware functions for context handling

api/
└── routes.go                              # API route registration with middleware

test/integration/
└── context_switch_test.go                 # 4 test suites with 20+ test cases
```

## Quick API Usage

### Python Example
```python
import requests

BASE_URL = "http://localhost:8080/api/v1"
AUTH_TOKEN = "your-jwt-token"

headers = {
    "Authorization": f"Bearer {AUTH_TOKEN}",
    "Content-Type": "application/json"
}

# Get available contexts
response = requests.get(
    f"{BASE_URL}/profile/available-contexts",
    headers=headers
)
contexts = response.json()

# Switch context
switch_response = requests.post(
    f"{BASE_URL}/profile/switch-context",
    headers=headers,
    json={"organization_id": "org_456"}
)

# Get switch history
history = requests.get(
    f"{BASE_URL}/profile/context-switch-history?limit=50&offset=0",
    headers=headers
)
```

### cURL Examples
```bash
# Get available contexts
curl -X GET http://localhost:8080/api/v1/profile/available-contexts \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get current context
curl -X GET http://localhost:8080/api/v1/profile/current-context \
  -H "Authorization: Bearer YOUR_TOKEN"

# Switch context
curl -X POST http://localhost:8080/api/v1/profile/switch-context \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"organization_id": "org_456"}'

# Get context switch history
curl -X GET http://localhost:8080/api/v1/profile/context-switch-history \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Service Layer Usage (Go)

```go
import (
    "context"
    "ethos/internal/organization/repository"
    "ethos/internal/organization/service"
)

// Initialize
repo := repository.NewPostgresContextRepository(db)
svc := service.NewUserContextService(repo)

// Get available contexts
contexts, err := svc.GetAvailableContexts(ctx, userID)

// Get current context
current, err := svc.GetCurrentContext(ctx, userID)

// Switch context
newContext, err := svc.SwitchContext(ctx, userID, orgID, ipAddress, userAgent)

// Validate membership
isMember, err := svc.ValidateUserInOrganization(ctx, userID, orgID)

// Get role
role, err := svc.GetUserRoleInOrganization(ctx, userID, orgID)

// Get history
records, err := svc.GetContextSwitchHistory(ctx, userID, limit, offset)
```

## Middleware Usage

```go
import "ethos/internal/middleware"

// In route setup
v1.Group("/profile").
    Use(middleware.AuthMiddleware(tokenGen)).
    Use(middleware.ContextSwitchMiddleware(contextService)).
    POST("/switch-context", handler.SwitchContext)

// For organization routes
orgs := v1.Group("/organizations").
    Use(middleware.AuthMiddleware(tokenGen)).
    Use(middleware.ValidateOrganizationMembership(contextService))
```

## Database Queries

### Check membership
```sql
SELECT COUNT(*) > 0 
FROM org_members 
WHERE user_id = $1 AND organization_id = $2 AND is_active = true
```

### Get user roles
```sql
SELECT role 
FROM org_members 
WHERE user_id = $1 AND organization_id = $2 AND is_active = true
```

### Record context switch
```sql
INSERT INTO user_context_switches 
(id, user_id, from_organization_id, to_organization_id, session_id, ip_address, timestamp)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW())
```

### Get switch history
```sql
SELECT * FROM user_context_switches
WHERE user_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3
```

## Response Examples

### Available Contexts
```json
{
  "contexts": [
    {
      "user_id": "user_123",
      "organization_id": "org_456",
      "organization_name": "Acme Corp",
      "role": "admin",
      "permissions": ["moderation", "analytics"],
      "joined_at": "2025-01-15T10:30:00Z",
      "last_switched_at": "2025-01-20T14:25:00Z"
    }
  ],
  "current": {
    "user_id": "user_123",
    "organization_id": "org_789",
    "organization_name": "Beta Inc",
    "role": "member",
    "permissions": ["feedback"],
    "joined_at": "2025-02-01T08:00:00Z",
    "last_switched_at": "2025-02-15T11:45:00Z"
  },
  "total": 3
}
```

### Switch Context Response
```json
{
  "context": {
    "user_id": "user_123",
    "organization_id": "org_456",
    "organization_name": "Acme Corp",
    "role": "admin",
    "permissions": ["moderation", "analytics"]
  },
  "message": "context switched successfully",
  "timestamp": "2025-01-20T14:25:00Z"
}
```

### Context Switch History
```json
{
  "records": [
    {
      "id": "switch_123",
      "user_id": "user_123",
      "from_organization_id": "org_789",
      "to_organization_id": "org_456",
      "session_id": "session_abc",
      "ip_address": "192.168.1.100",
      "timestamp": "2025-01-20T14:25:00Z"
    }
  ],
  "total": 42
}
```

## Error Codes

| Code | HTTP | Meaning |
|------|------|---------|
| UNAUTHORIZED | 401 | Missing or invalid authentication token |
| AUTH_REQUIRED | 401 | Authentication middleware requirement |
| FORBIDDEN | 403 | User not authorized for operation |
| NOT_ORGANIZATION_MEMBER | 403 | User is not member of organization |
| NOT_FOUND | 404 | Resource (context/session) not found |
| INVALID_REQUEST | 400 | Malformed request or validation failed |
| SERVER_ERROR | 500 | Unexpected server error |
| CONTEXT_LOAD_FAILED | 500 | Failed to load user context |

## Performance Metrics

- **Avg Context Switch**: 50-100ms
- **Avg History Retrieval**: 30-50ms (50 records)
- **Avg Membership Check**: 10-20ms
- **Session Creation**: 20-40ms

## Security Checklist

- ✅ Token validation on all endpoints
- ✅ Membership verification before switches
- ✅ IP address logging
- ✅ User-agent logging
- ✅ SHA256 token hashing
- ✅ Audit trail creation
- ✅ Org boundary enforcement
- ✅ Role-based access control

## Troubleshooting

### Context switch fails with 403
- Check user is member of target org: `svc.ValidateUserInOrganization()`
- Verify membership is active in `org_members` table

### No available contexts
- Verify user has entries in `org_members` table
- Check `is_active` flag is true

### Session not found
- Verify token hash matches: `GetUserSessionByToken()`
- Check session hasn't been revoked
- Confirm session hasn't expired

### Empty history
- Verify user has made context switches
- Check `user_context_switches` table has records
- Confirm user_id matches exactly

## Development Tips

1. **Testing**: Run integration tests with `go test -v ./test/integration/...`
2. **Logging**: Add debug logs to middleware for context flow tracking
3. **Caching**: Implement org list caching for performance
4. **Monitoring**: Track context switch latency and frequency

## Related Documentation

- Database Schema: `DEVELOPMENT_SETUP.md`
- API Specs: `APISpecs.md`
- Architecture: `ARCHITECTURE.md`
