package repository

import (
	"context"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertRepository interface {
	CreateAlert(ctx context.Context, alert *model.AssetAlert) error
	GetAlertsByAssetID(ctx context.Context, userID uuid.UUID, assetID uuid.UUID) ([]model.AssetAlert, error)
	UpdateAlert(ctx context.Context, alert *model.AssetAlert) error
	DeleteAlert(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type alertRepository struct {
	db *pgxpool.Pool
}

func NewAlertRepository(db *pgxpool.Pool) AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) CreateAlert(ctx context.Context, alert *model.AssetAlert) error {
	query := `
		INSERT INTO asset_alerts (
			user_id, asset_id, alert_type, target_value, is_active, is_triggered, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		alert.UserID, alert.AssetID, alert.AlertType, alert.TargetValue,
		alert.IsActive, alert.IsTriggered, alert.Notes,
	).Scan(&alert.ID, &alert.CreatedAt, &alert.UpdatedAt)

	return err
}

func (r *alertRepository) GetAlertsByAssetID(ctx context.Context, userID uuid.UUID, assetID uuid.UUID) ([]model.AssetAlert, error) {
	query := `
		SELECT id, user_id, asset_id, alert_type, target_value, is_active, is_triggered, notes, created_at, updated_at
		FROM asset_alerts
		WHERE user_id = $1 AND asset_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []model.AssetAlert
	for rows.Next() {
		var a model.AssetAlert
		err := rows.Scan(
			&a.ID, &a.UserID, &a.AssetID, &a.AlertType, &a.TargetValue,
			&a.IsActive, &a.IsTriggered, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *alertRepository) UpdateAlert(ctx context.Context, alert *model.AssetAlert) error {
	query := `
		UPDATE asset_alerts
		SET alert_type = $1, target_value = $2, is_active = $3, is_triggered = $4, notes = $5, updated_at = NOW()
		WHERE id = $6 AND user_id = $7
	`
	_, err := r.db.Exec(ctx, query,
		alert.AlertType, alert.TargetValue, alert.IsActive, alert.IsTriggered, alert.Notes,
		alert.ID, alert.UserID,
	)
	return err
}

func (r *alertRepository) DeleteAlert(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM asset_alerts WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}
