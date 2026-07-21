package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	GetByUserID(ctx context.Context, userID string) ([]*model.Notification, error)
	MarkAsRead(ctx context.Context, notifID string, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}

type notificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notif *model.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, message, action_link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_read, created_at
	`
	return r.db.QueryRow(ctx, query,
		notif.UserID, notif.Type, notif.Title, notif.Message, notif.ActionLink,
	).Scan(&notif.ID, &notif.IsRead, &notif.CreatedAt)
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, action_link, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []*model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message,
			&n.IsRead, &n.ActionLink, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifs = append(notifs, &n)
	}
	return notifs, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notifID string, userID string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, notifID, userID)
	return err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
