package service

import (
	"context"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, userID, notifType, title, message string, actionLink *string) error
	GetNotifications(ctx context.Context, userID string) ([]*model.Notification, error)
	MarkAsRead(ctx context.Context, notifID, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) CreateNotification(ctx context.Context, userID, notifType, title, message string, actionLink *string) error {
	notif := &model.Notification{
		UserID:     userID,
		Type:       notifType,
		Title:      title,
		Message:    message,
		ActionLink: actionLink,
	}
	return s.repo.Create(ctx, notif)
}

func (s *notificationService) GetNotifications(ctx context.Context, userID string) ([]*model.Notification, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, notifID, userID string) error {
	return s.repo.MarkAsRead(ctx, notifID, userID)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}
