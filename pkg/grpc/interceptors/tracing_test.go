package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestTracingUnaryClientInterceptor(t *testing.T) {
	tracer := otel.Tracer("test")
	interceptor := TracingUnaryClientInterceptor()

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Create a test invoker that checks context propagation
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		// Verify context has span
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.True(t, ok)
		assert.NotNil(t, md)
		return nil
	}

	// Test interceptor
	err := interceptor(ctx, "TestMethod", nil, nil, nil, invoker)
	assert.NoError(t, err)
}

func TestTracingStreamClientInterceptor(t *testing.T) {
	tracer := otel.Tracer("test")
	interceptor := TracingStreamClientInterceptor()

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Create a test streamer
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// Verify context has span
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.True(t, ok)
		assert.NotNil(t, md)
		return nil, nil
	}

	// Test interceptor
	stream, err := interceptor(ctx, nil, nil, "TestMethod", streamer)
	assert.NoError(t, err)
	assert.Nil(t, stream) // Stream is nil in test
}

