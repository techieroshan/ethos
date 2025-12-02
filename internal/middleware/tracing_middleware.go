package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer  = otel.Tracer("ethos-api")
	meter   = otel.Meter("ethos-api")
	requestCounter metric.Int64Counter
	requestDuration metric.Float64Histogram
)

func init() {
	var err error
	requestCounter, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		// Log error but don't fail
	}

	requestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
	)
	if err != nil {
		// Log error but don't fail
	}
}

// TracingMiddleware adds OpenTelemetry tracing to requests
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		ctx, span := tracer.Start(
			c.Request.Context(),
			c.FullPath(),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Set request attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.route", c.FullPath()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
		)

		// Update context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		statusCode := c.Writer.Status()

		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Float64("http.duration", duration),
		)

		if requestCounter != nil {
			requestCounter.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("http.method", c.Request.Method),
					attribute.String("http.route", c.FullPath()),
					attribute.Int("http.status_code", statusCode),
				),
			)
		}

		if requestDuration != nil {
			requestDuration.Record(ctx, duration,
				metric.WithAttributes(
					attribute.String("http.method", c.Request.Method),
					attribute.String("http.route", c.FullPath()),
					attribute.Int("http.status_code", statusCode),
				),
			)
		}

		// Set span status based on HTTP status code
		if statusCode >= 400 {
			span.SetStatus(codes.Error, c.Errors.String())
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}
}

