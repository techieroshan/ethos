package interceptors

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TracingUnaryClientInterceptor returns a gRPC unary client interceptor for tracing
func TracingUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		tracer := otel.Tracer("grpc-client")
		ctx, span := tracer.Start(ctx, method,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		// Add attributes
		span.SetAttributes(
			attribute.String("rpc.method", method),
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", extractServiceName(method)),
		)

		// Inject trace context into metadata
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			// Copy existing metadata
			md = md.Copy()
		}

		// Create a carrier for metadata
		carrier := make(metadata.MD)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(carrier))
		
		// Merge carrier into metadata
		for k, v := range carrier {
			md[k] = v
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Invoke the RPC
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Record result
		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				span.SetStatus(codes.Error, s.Message())
				span.SetAttributes(attribute.String("rpc.status_code", s.Code().String()))
			} else {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return err
	}
}

// TracingStreamClientInterceptor returns a gRPC stream client interceptor for tracing
func TracingStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		tracer := otel.Tracer("grpc-client")
		ctx, span := tracer.Start(ctx, method,
			trace.WithSpanKind(trace.SpanKindClient),
		)

		// Add attributes
		span.SetAttributes(
			attribute.String("rpc.method", method),
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", extractServiceName(method)),
		)

		// Inject trace context into metadata
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			// Copy existing metadata
			md = md.Copy()
		}

		// Create a carrier for metadata
		carrier := make(metadata.MD)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(carrier))
		
		// Merge carrier into metadata
		for k, v := range carrier {
			md[k] = v
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Create stream
		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			return nil, err
		}

		// Wrap stream to close span when done
		wrappedStream := &tracedClientStream{
			ClientStream: stream,
			span:         span,
		}

		return wrappedStream, nil
	}
}

// tracedClientStream wraps a ClientStream to close the span when done
type tracedClientStream struct {
	grpc.ClientStream
	span trace.Span
}

func (s *tracedClientStream) CloseSend() error {
	err := s.ClientStream.CloseSend()
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (s *tracedClientStream) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
		s.span.End()
	}
	return err
}

func (s *tracedClientStream) SendMsg(m interface{}) error {
	err := s.ClientStream.SendMsg(m)
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// extractServiceName extracts the service name from a gRPC method name
// Example: "/feedback.FeedbackService/GetFeed" -> "feedback.FeedbackService"
func extractServiceName(method string) string {
	if len(method) == 0 {
		return "unknown"
	}
	// Remove leading slash
	if method[0] == '/' {
		method = method[1:]
	}
	// Extract service name (everything before the last slash)
	for i := len(method) - 1; i >= 0; i-- {
		if method[i] == '/' {
			return method[:i]
		}
	}
	return method
}
