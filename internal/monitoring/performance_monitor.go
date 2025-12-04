package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1model "github.com/prometheus/common/model"
)

// PerformanceMonitor handles application performance monitoring
type PerformanceMonitor struct {
	// Prometheus metrics
	requestDuration *prometheus.HistogramVec
	requestCount    *prometheus.CounterVec
	activeConnections *prometheus.GaugeVec
	cacheHitRatio   *prometheus.GaugeVec
	dbQueryDuration *prometheus.HistogramVec
	rateLimitHits   *prometheus.CounterVec

	// Internal metrics storage
	metricsMutex sync.RWMutex
	requestMetrics map[string]*RequestMetrics
}

// RequestMetrics holds metrics for a specific endpoint
type RequestMetrics struct {
	TotalRequests int64
	TotalDuration time.Duration
	AvgDuration   time.Duration
	LastRequest   time.Time
	ErrorCount    int64
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	pm := &PerformanceMonitor{
		requestMetrics: make(map[string]*RequestMetrics),
	}

	// Initialize Prometheus metrics
	pm.requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	pm.requestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	pm.activeConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
		[]string{"type"}, // api, db, cache
	)

	pm.cacheHitRatio = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_hit_ratio",
			Help: "Cache hit ratio percentage",
		},
		[]string{"cache_type"},
	)

	pm.dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"query_type", "table"},
	)

	pm.rateLimitHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"endpoint", "client_ip"},
	)

	return pm
}

// RecordHTTPRequest records an HTTP request metric
func (pm *PerformanceMonitor) RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	pm.requestDuration.WithLabelValues(method, endpoint, status).Observe(duration.Seconds())
	pm.requestCount.WithLabelValues(method, endpoint, status).Inc()

	// Update internal metrics
	pm.metricsMutex.Lock()
	defer pm.metricsMutex.Unlock()

	key := fmt.Sprintf("%s:%s", method, endpoint)
	if pm.requestMetrics[key] == nil {
		pm.requestMetrics[key] = &RequestMetrics{}
	}

	metrics := pm.requestMetrics[key]
	metrics.TotalRequests++
	metrics.TotalDuration += duration
	metrics.AvgDuration = metrics.TotalDuration / time.Duration(metrics.TotalRequests)
	metrics.LastRequest = time.Now()

	if status[0] == '5' || status[0] == '4' {
		metrics.ErrorCount++
	}
}

// RecordDBQuery records a database query metric
func (pm *PerformanceMonitor) RecordDBQuery(queryType, table string, duration time.Duration) {
	pm.dbQueryDuration.WithLabelValues(queryType, table).Observe(duration.Seconds())
}

// UpdateActiveConnections updates the active connections gauge
func (pm *PerformanceMonitor) UpdateActiveConnections(connType string, count float64) {
	pm.activeConnections.WithLabelValues(connType).Set(count)
}

// UpdateCacheHitRatio updates the cache hit ratio gauge
func (pm *PerformanceMonitor) UpdateCacheHitRatio(cacheType string, ratio float64) {
	pm.cacheHitRatio.WithLabelValues(cacheType).Set(ratio)
}

// RecordRateLimitHit records a rate limit hit
func (pm *PerformanceMonitor) RecordRateLimitHit(endpoint, clientIP string) {
	pm.rateLimitHits.WithLabelValues(endpoint, clientIP).Inc()
}

// GetRequestMetrics returns current request metrics
func (pm *PerformanceMonitor) GetRequestMetrics() map[string]*RequestMetrics {
	pm.metricsMutex.RLock()
	defer pm.metricsMutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]*RequestMetrics)
	for k, v := range pm.requestMetrics {
		result[k] = &RequestMetrics{
			TotalRequests: v.TotalRequests,
			TotalDuration: v.TotalDuration,
			AvgDuration:   v.AvgDuration,
			LastRequest:   v.LastRequest,
			ErrorCount:    v.ErrorCount,
		}
	}
	return result
}

// GetSlowEndpoints returns endpoints with average response time above threshold
func (pm *PerformanceMonitor) GetSlowEndpoints(threshold time.Duration) []string {
	pm.metricsMutex.RLock()
	defer pm.metricsMutex.RUnlock()

	var slowEndpoints []string
	for endpoint, metrics := range pm.requestMetrics {
		if metrics.AvgDuration > threshold {
			slowEndpoints = append(slowEndpoints, endpoint)
		}
	}
	return slowEndpoints
}

// GetHighErrorRateEndpoints returns endpoints with error rate above threshold
func (pm *PerformanceMonitor) GetHighErrorRateEndpoints(errorRateThreshold float64) []string {
	pm.metricsMutex.RLock()
	defer pm.metricsMutex.RUnlock()

	var highErrorEndpoints []string
	for endpoint, metrics := range pm.requestMetrics {
		if metrics.TotalRequests > 0 {
			errorRate := float64(metrics.ErrorCount) / float64(metrics.TotalRequests)
			if errorRate > errorRateThreshold {
				highErrorEndpoints = append(highErrorEndpoints, endpoint)
			}
		}
	}
	return highErrorEndpoints
}

// HealthChecker interface for health checking components
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
	Name() string
}

// ComponentHealth represents the health status of a component
type ComponentHealth struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // "healthy", "unhealthy", "degraded"
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	ResponseTime time.Duration `json:"response_time,omitempty"`
}

// HealthMonitor monitors the health of system components
type HealthMonitor struct {
	checkers []HealthChecker
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor() *HealthMonitor {
	return &HealthMonitor{
		checkers: make([]HealthChecker, 0),
	}
}

// AddChecker adds a health checker to the monitor
func (hm *HealthMonitor) AddChecker(checker HealthChecker) {
	hm.checkers = append(hm.checkers, checker)
}

// CheckHealth performs health checks on all registered components
func (hm *HealthMonitor) CheckHealth(ctx context.Context) []ComponentHealth {
	results := make([]ComponentHealth, 0, len(hm.checkers))

	for _, checker := range hm.checkers {
		start := time.Now()
		err := checker.HealthCheck(ctx)
		duration := time.Since(start)

		health := ComponentHealth{
			Name:         checker.Name(),
			Timestamp:    time.Now(),
			ResponseTime: duration,
		}

		if err != nil {
			health.Status = "unhealthy"
			health.Message = err.Error()
		} else {
			health.Status = "healthy"
		}

		results = append(results, health)
	}

	return results
}

// GetOverallHealth returns the overall system health status
func (hm *HealthMonitor) GetOverallHealth(ctx context.Context) ComponentHealth {
	healths := hm.CheckHealth(ctx)

	overall := ComponentHealth{
		Name:      "system",
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	for _, health := range healths {
		if health.Status == "unhealthy" {
			overall.Status = "unhealthy"
			overall.Message = fmt.Sprintf("Component %s is unhealthy: %s", health.Name, health.Message)
			break
		}
		if health.Status == "degraded" && overall.Status == "healthy" {
			overall.Status = "degraded"
			overall.Message = fmt.Sprintf("Component %s is degraded", health.Name)
		}
	}

	return overall
}

// PrometheusClient provides access to Prometheus metrics
type PrometheusClient struct {
	client api.Client
	v1api  v1.API
}

// NewPrometheusClient creates a new Prometheus client
func NewPrometheusClient(prometheusURL string) (*PrometheusClient, error) {
	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}

	v1api := v1.NewAPI(client)

	return &PrometheusClient{
		client: client,
		v1api:  v1api,
	}, nil
}

// QueryInstant queries Prometheus for instant metrics
func (pc *PrometheusClient) QueryInstant(ctx context.Context, query string) (float64, error) {
	result, warnings, err := pc.v1api.Query(ctx, query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to query Prometheus: %w", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("Prometheus warnings: %v\n", warnings)
	}

	// Extract scalar value
	if result.Type().String() == "scalar" {
		if scalar, ok := result.(*v1model.Scalar); ok {
			return float64(scalar.Value), nil
		}
	}

	return 0, fmt.Errorf("unexpected result type: %s", result.Type().String())
}

// QueryRange queries Prometheus for range metrics
func (pc *PrometheusClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]*v1model.SamplePair, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	result, warnings, err := pc.v1api.QueryRange(ctx, query, r)
	if err != nil {
		return nil, fmt.Errorf("failed to query range: %w", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("Prometheus warnings: %v\n", warnings)
	}

	// Extract matrix data
	if result.Type().String() == "matrix" {
		if matrix, ok := result.(v1model.Matrix); ok && len(matrix) > 0 {
			// Convert []model.SamplePair to []*model.SamplePair
			values := make([]*v1model.SamplePair, len(matrix[0].Values))
			for i, v := range matrix[0].Values {
				values[i] = &v
			}
			return values, nil
		}
	}

	return nil, fmt.Errorf("no data returned or unexpected result type")
}
