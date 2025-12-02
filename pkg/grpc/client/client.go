package client

import (
	"sync"
	"time"

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
			)
			if err != nil {
				// Return nil if connection fails (will be handled by retry logic)
				return nil
			}
			m.feedbackConn = conn
		}
		// TODO: Create actual client after proto generation
		// m.feedbackClient = feedbackpb.NewFeedbackServiceClient(m.feedbackConn)
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
			)
			if err != nil {
				return nil
			}
			m.dashboardConn = conn
		}
		// TODO: Create actual client after proto generation
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
			)
			if err != nil {
				return nil
			}
			m.notificationsConn = conn
		}
		// TODO: Create actual client after proto generation
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
			)
			if err != nil {
				return nil
			}
			m.peopleConn = conn
		}
		// TODO: Create actual client after proto generation
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

