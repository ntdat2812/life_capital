package model

import (
	"time"

	"github.com/google/uuid"
)

type AssetAlert struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	AssetID     uuid.UUID `json:"asset_id" db:"asset_id"`
	AlertType   string    `json:"alert_type" db:"alert_type"`       // 'take_profit', 'stop_loss', 'stop_accumulating', 'custom'
	TargetValue float64   `json:"target_value" db:"target_value"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsTriggered bool      `json:"is_triggered" db:"is_triggered"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAlertRequest struct {
	AlertType   string  `json:"alert_type" validate:"required"`
	TargetValue float64 `json:"target_value" validate:"required,gt=0"`
	Notes       string  `json:"notes,omitempty"`
}

type UpdateAlertRequest struct {
	AlertType   string  `json:"alert_type" validate:"required"`
	TargetValue float64 `json:"target_value" validate:"required,gt=0"`
	IsActive    bool    `json:"is_active"`
	IsTriggered bool    `json:"is_triggered"`
	Notes       string  `json:"notes,omitempty"`
}

