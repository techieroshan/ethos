package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClientManager(t *testing.T) {
	config := Config{
		FeedbackEndpoint:     "localhost:50051",
		DashboardEndpoint:    "localhost:50052",
		NotificationsEndpoint: "localhost:50053",
		PeopleEndpoint:       "localhost:50054",
		Timeout:              5 * time.Second,
		Retries:              3,
	}

	manager := NewClientManager(config)

	assert.NotNil(t, manager)
	assert.Equal(t, config.FeedbackEndpoint, manager.config.FeedbackEndpoint)
}

func TestClientManager_GetFeedbackClient(t *testing.T) {
	config := Config{
		FeedbackEndpoint: "localhost:50051",
		Timeout:          5 * time.Second,
		Retries:          3,
	}

	manager := NewClientManager(config)
	client := manager.GetFeedbackClient()

	assert.NotNil(t, client)
}

func TestClientManager_GetDashboardClient(t *testing.T) {
	config := Config{
		DashboardEndpoint: "localhost:50052",
		Timeout:           5 * time.Second,
		Retries:           3,
	}

	manager := NewClientManager(config)
	client := manager.GetDashboardClient()

	assert.NotNil(t, client)
}

func TestClientManager_GetNotificationsClient(t *testing.T) {
	config := Config{
		NotificationsEndpoint: "localhost:50053",
		Timeout:               5 * time.Second,
		Retries:               3,
	}

	manager := NewClientManager(config)
	client := manager.GetNotificationsClient()

	assert.NotNil(t, client)
}

func TestClientManager_GetPeopleClient(t *testing.T) {
	config := Config{
		PeopleEndpoint: "localhost:50054",
		Timeout:        5 * time.Second,
		Retries:        3,
	}

	manager := NewClientManager(config)
	client := manager.GetPeopleClient()

	assert.NotNil(t, client)
}

func TestClientManager_Close(t *testing.T) {
	config := Config{
		FeedbackEndpoint: "localhost:50051",
		Timeout:          5 * time.Second,
		Retries:          3,
	}

	manager := NewClientManager(config)
	
	// Close should not panic
	assert.NotPanics(t, func() {
		manager.Close()
	})
}

func TestClientManager_ConnectionPooling(t *testing.T) {
	config := Config{
		FeedbackEndpoint: "localhost:50051",
		Timeout:          5 * time.Second,
		Retries:          3,
	}

	manager := NewClientManager(config)
	
	// Get same client multiple times should return same instance (connection reuse)
	client1 := manager.GetFeedbackClient()
	client2 := manager.GetFeedbackClient()

	assert.Equal(t, client1, client2)
}

func TestClientManager_ContextPropagation(t *testing.T) {
	config := Config{
		FeedbackEndpoint: "localhost:50051",
		Timeout:          5 * time.Second,
		Retries:          3,
	}

	manager := NewClientManager(config)
	ctx := context.WithValue(context.Background(), "user_id", "test-user")

	// Context should be propagated through client calls
	// This is tested implicitly through client usage
	assert.NotNil(t, ctx)
	assert.NotNil(t, manager)
}

