package client

import (
	"sync"
	"time"

	"ethos/pkg/grpc/interceptors"
	feedbackpb "ethos/api/proto/feedback"
	dashboardpb "ethos/api/proto/dashboard"
	notificationpb "ethos/api/proto/notifications"
	peoplepb "ethos/api/proto/people"
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
	feedbackClient      feedbackpb.FeedbackServiceClient
	dashboardClient     dashboardpb.DashboardServiceClient
	notificationsClient notificationpb.NotificationServiceClient
	peopleClient        peoplepb.PeopleServiceClient
	mu                  sync.RWMutex
}

// NewClientManager creates a new gRPC client manager
func NewClientManager(config Config) *ClientManager {
	return &ClientManager{
		config: config,
	}
}

// GetFeedbackClient returns or creates the feedback gRPC client
func (m *ClientManager) GetFeedbackClient() feedbackpb.FeedbackServiceClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.feedbackClient == nil {
		// Create connection if needed
		if m.feedbackConn == nil {
			conn, err := grpc.NewClient(
				m.config.FeedbackEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				// Return nil if connection fails (caller should handle)
				return nil
			}
			m.feedbackConn = conn
		}
		// Create actual client from generated proto
		m.feedbackClient = feedbackpb.NewFeedbackServiceClient(m.feedbackConn)
	}
	return m.feedbackClient
}

// GetDashboardClient returns or creates the dashboard gRPC client
func (m *ClientManager) GetDashboardClient() dashboardpb.DashboardServiceClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.dashboardClient == nil {
		if m.dashboardConn == nil {
			conn, err := grpc.NewClient(
				m.config.DashboardEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				return nil
			}
			m.dashboardConn = conn
		}
		m.dashboardClient = dashboardpb.NewDashboardServiceClient(m.dashboardConn)
	}
	return m.dashboardClient
}

// GetNotificationsClient returns or creates the notifications gRPC client
func (m *ClientManager) GetNotificationsClient() notificationpb.NotificationServiceClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.notificationsClient == nil {
		if m.notificationsConn == nil {
			conn, err := grpc.NewClient(
				m.config.NotificationsEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				return nil
			}
			m.notificationsConn = conn
		}
		m.notificationsClient = notificationpb.NewNotificationServiceClient(m.notificationsConn)
	}
	return m.notificationsClient
}

// GetPeopleClient returns or creates the people gRPC client
func (m *ClientManager) GetPeopleClient() peoplepb.PeopleServiceClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.peopleClient == nil {
		if m.peopleConn == nil {
			conn, err := grpc.NewClient(
				m.config.PeopleEndpoint,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(interceptors.TracingUnaryClientInterceptor()),
				grpc.WithStreamInterceptor(interceptors.TracingStreamClientInterceptor()),
			)
			if err != nil {
				return nil
			}
			m.peopleConn = conn
		}
		m.peopleClient = peoplepb.NewPeopleServiceClient(m.peopleConn)
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

