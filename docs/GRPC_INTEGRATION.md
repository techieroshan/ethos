# gRPC Integration Guide

## Overview

This document describes the gRPC integration strategy for the Ethos BFF service. gRPC is used selectively for performance-critical operations between the BFF and Core Backend services.

## Architecture

The BFF uses a **hybrid approach**:
- **REST** for simple CRUD operations and low-frequency endpoints
- **gRPC** for high-frequency read operations, complex queries, and aggregations

## Protocol Selection

Protocol selection is **configuration-driven** via environment variables:

```bash
GRPC_ENABLED=true
GRPC_FEEDBACK_ENDPOINT=localhost:50051
GRPC_DASHBOARD_ENDPOINT=localhost:50052
GRPC_NOTIFICATIONS_ENDPOINT=localhost:50053
GRPC_PEOPLE_ENDPOINT=localhost:50054
```

When `GRPC_ENABLED=false` or endpoints are not configured, the BFF falls back to REST.

## Implementation Status

### Completed

1. ✅ Protocol Buffer definitions (`.proto` files)
2. ✅ gRPC client manager with connection pooling
3. ✅ OpenTelemetry tracing interceptors
4. ✅ Configuration support
5. ✅ Email services (Checker, Emailit, Mailpit)

### Pending (Requires Proto Code Generation)

1. ⏳ Generate Go code from `.proto` files
2. ⏳ Create service-specific gRPC client wrappers
3. ⏳ Implement REST/gRPC adapter pattern
4. ⏳ Update BFF services to use gRPC clients
5. ⏳ Integration tests

## Protocol Buffer Code Generation

### Prerequisites

Install Protocol Buffer compiler and Go plugins:

```bash
# Install protoc (macOS)
brew install protobuf

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Generate Code

Run the Makefile target:

```bash
make proto
```

Or manually:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  api/proto/**/*.proto
```

### Generated Files

After generation, you'll have:
- `api/proto/feedback/feedback.pb.go` - Feedback messages
- `api/proto/feedback/feedback_grpc.pb.go` - Feedback service client
- Similar files for dashboard, notifications, people

## Service-Specific Clients

After proto generation, update `pkg/grpc/client/client.go` to use generated clients:

```go
import feedbackpb "ethos/api/proto/feedback"

func (m *ClientManager) GetFeedbackClient() feedbackpb.FeedbackServiceClient {
    // Replace stub with actual client
    return feedbackpb.NewFeedbackServiceClient(m.feedbackConn)
}
```

## REST/gRPC Adapter Pattern

Services use an adapter pattern to switch between REST and gRPC:

```go
type FeedbackServiceClient interface {
    GetFeed(ctx context.Context, limit, offset int) ([]*FeedbackItem, int, error)
}

type RESTFeedbackClient struct {
    // REST implementation
}

type GRPCFeedbackClient struct {
    client feedbackpb.FeedbackServiceClient
}

// Service uses interface, switches based on config
```

## Performance Expectations

- **Latency**: 20-40% reduction for gRPC vs REST
- **Payload Size**: 30-50% reduction (binary vs JSON)
- **Throughput**: Higher with HTTP/2 multiplexing

## Testing

### Local Development

1. Start Mailpit for email testing: `docker run -d -p 1025:1025 -p 8025:8025 axllent/mailpit`
2. Set `MAILPIT_ENABLED=true` in `.env`
3. Use REST for all operations (easier debugging)

### Integration Testing

1. Start Core Backend gRPC servers (or use mocks)
2. Set `GRPC_ENABLED=true` in test environment
3. Run integration tests

## Troubleshooting

### Connection Errors

- Verify gRPC server is running
- Check endpoint configuration
- Review connection timeout settings

### Proto Generation Issues

- Ensure `protoc` is in PATH
- Verify Go plugins are installed
- Check `.proto` file syntax

### Tracing Not Working

- Verify OpenTelemetry is enabled
- Check interceptor registration
- Review Jaeger/collector configuration

## Next Steps

1. Generate proto code: `make proto`
2. Update client manager to use generated clients
3. Create service adapters (REST/gRPC)
4. Update BFF services with feature flags
5. Add integration tests
6. Performance benchmarking

