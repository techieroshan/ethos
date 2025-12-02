package client

import (
	"strings"

	"ethos/internal/config"
	feedbackService "ethos/internal/feedback/service"
	dashboardService "ethos/internal/dashboard/service"
	notificationService "ethos/internal/notifications/service"
	peopleService "ethos/internal/people/service"
	feedbackRepo "ethos/internal/feedback/repository"
	dashboardRepo "ethos/internal/dashboard/repository"
	notificationRepo "ethos/internal/notifications/repository"
	peopleRepo "ethos/internal/people/repository"
)

// CreateFeedbackService creates a feedback service based on protocol configuration
func CreateFeedbackService(cfg *config.Config, grpcManager *ClientManager, repo feedbackRepo.Repository) feedbackService.Service {
	if strings.ToLower(cfg.GRPC.FeedbackProtocol) == "grpc" && cfg.GRPC.Enabled {
		grpcClient := grpcManager.GetFeedbackClient()
		if grpcClient != nil {
			client := feedbackService.NewGRPCFeedbackClient(grpcClient)
			return feedbackService.NewFeedbackServiceWithClient(client, repo)
		}
	}
	// Default to REST
	return feedbackService.NewFeedbackService(repo)
}

// CreateDashboardService creates a dashboard service based on protocol configuration
func CreateDashboardService(cfg *config.Config, grpcManager *ClientManager, repo dashboardRepo.Repository) dashboardService.Service {
	if strings.ToLower(cfg.GRPC.DashboardProtocol) == "grpc" && cfg.GRPC.Enabled {
		grpcClient := grpcManager.GetDashboardClient()
		if grpcClient != nil {
			client := dashboardService.NewGRPCDashboardClient(grpcClient)
			return dashboardService.NewDashboardServiceWithClient(client)
		}
	}
	// Default to REST
	return dashboardService.NewDashboardService(repo)
}

// CreateNotificationService creates a notification service based on protocol configuration
func CreateNotificationService(cfg *config.Config, grpcManager *ClientManager, repo notificationRepo.Repository) notificationService.Service {
	if strings.ToLower(cfg.GRPC.NotificationsProtocol) == "grpc" && cfg.GRPC.Enabled {
		grpcClient := grpcManager.GetNotificationsClient()
		if grpcClient != nil {
			client := notificationService.NewGRPCNotificationClient(grpcClient)
			return notificationService.NewNotificationServiceWithClient(client, repo)
		}
	}
	// Default to REST
	return notificationService.NewNotificationService(repo)
}

// CreatePeopleService creates a people service based on protocol configuration
func CreatePeopleService(cfg *config.Config, grpcManager *ClientManager, repo peopleRepo.Repository) peopleService.Service {
	if strings.ToLower(cfg.GRPC.PeopleProtocol) == "grpc" && cfg.GRPC.Enabled {
		grpcClient := grpcManager.GetPeopleClient()
		if grpcClient != nil {
			client := peopleService.NewGRPCPeopleClient(grpcClient)
			return peopleService.NewPeopleServiceWithClient(client)
		}
	}
	// Default to REST
	return peopleService.NewPeopleService(repo)
}

