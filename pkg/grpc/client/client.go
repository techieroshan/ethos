package client

import (
	"sync"
	"time"

	"ethos/pkg/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds gRPC client configuration
type Config struct {
	FeedbackEndpoint     string
	DashboardEndpoint    string
	NotificationsEndpoint string
	PeopleEndpoint       string
	Timeout              time.Duration
	Retries              int
}

// ClientManager manages gRPC client connections
type ClientManager struct {
	config              Config
	feedbackConn        *grpc.ClientConn
	dashboardConn       *grpc.ClientConn
	notificationsConn   *grpc.ClientConn
	peopleConn          *grpc.ClientConn
	feedbackClient      interface{} // Will be typed after proto generation
	dashboardClient     interface{} // Will be typed after proto generation
	notificationsClient interface{} // Will be typed after proto generation
	peopleClient        interface{} // Will be typed after proto generation
	mu                  sync.RWMutex
}

// NewClientManager creates a new gRPC client manager
func NewClientManager(config Config) *ClientManager {
	return &ClientManager{
		config: config,
	}
}

// GetFeedbackClient returns or creates the feedback gRPC client
func (m *ClientManager) GetFeedbackClient() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.feedbackClient == nil {
		// Create connection if needed
		if m.feedbackConn == nil {
			conn, err := grpc.NewClient(
				m.config.FeedbackEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithTimeout(m.config.Timeout),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				// Return stub client for testing (will be replaced after proto generation)
				m.feedbackClient = &stubFeedbackClient{}
				return m.feedbackClient
			}
			m.feedbackConn = conn
		}
		// TODO: Create actual client after proto generation
		// m.feedbackClient = feedbackpb.NewFeedbackServiceClient(m.feedbackConn)
		// For now, return stub client
		if m.feedbackClient == nil {
			m.feedbackClient = &stubFeedbackClient{}
		}
	}
	return m.feedbackClient
}

// GetDashboardClient returns or creates the dashboard gRPC client
func (m *ClientManager) GetDashboardClient() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.dashboardClient == nil {
		if m.dashboardConn == nil {
			conn, err := grpc.NewClient(
				m.config.DashboardEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithTimeout(m.config.Timeout),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				m.dashboardClient = &stubDashboardClient{}
				return m.dashboardClient
			}
			m.dashboardConn = conn
		}
		// TODO: Create actual client after proto generation
		if m.dashboardClient == nil {
			m.dashboardClient = &stubDashboardClient{}
		}
	}
	return m.dashboardClient
}

// GetNotificationsClient returns or creates the notifications gRPC client
func (m *ClientManager) GetNotificationsClient() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.notificationsClient == nil {
		if m.notificationsConn == nil {
			conn, err := grpc.NewClient(
				m.config.NotificationsEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithTimeout(m.config.Timeout),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				m.notificationsClient = &stubNotificationsClient{}
				return m.notificationsClient
			}
			m.notificationsConn = conn
		}
		// TODO: Create actual client after proto generation
		if m.notificationsClient == nil {
			m.notificationsClient = &stubNotificationsClient{}
		}
	}
	return m.notificationsClient
}

// GetPeopleClient returns or creates the people gRPC client
func (m *ClientManager) GetPeopleClient() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.peopleClient == nil {
		if m.peopleConn == nil {
			conn, err := grpc.NewClient(
				m.config.PeopleEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithTimeout(m.config.Timeout),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				m.peopleClient = &stubPeopleClient{}
				return m.peopleClient
			}
			m.peopleConn = conn
		}
		// TODO: Create actual client after proto generation
		if m.peopleClient == nil {
			m.peopleClient = &stubPeopleClient{}
		}
	}
	return m.peopleClient
}

// Close closes all gRPC connections
func (m *ClientManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	if m.feedbackConn != nil {
		if err := m.feedbackConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if m.dashboardConn != nil {
		if err := m.dashboardConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if m.notificationsConn != nil {
		if err := m.notificationsConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if m.peopleConn != nil {
		if err := m.peopleConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0] // Return first error
	}
	return nil
}

// stubFeedbackClient is a stub client for testing (will be replaced after proto generation)
type stubFeedbackClient struct{}

// stubDashboardClient is a stub client for testing
type stubDashboardClient struct{}

// stubNotificationsClient is a stub client for testing
type stubNotificationsClient struct{}

// stubPeopleClient is a stub client for testing
type stubPeopleClient struct{}

