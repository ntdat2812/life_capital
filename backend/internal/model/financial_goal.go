package model

import (
	"time"

	"github.com/google/uuid"
)

type FinancialGoal struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Name         string     `json:"name" db:"name"`
	TargetAmount float64    `json:"target_amount" db:"target_amount"`
	Priority     int        `json:"priority" db:"priority"`
	TargetDate   *time.Time `json:"target_date" db:"target_date"`
	Status       string     `json:"status" db:"status"` // active, completed
	CreatedAt    time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at" db:"updated_at"`
	Allocations  []GoalAssetAllocation `json:"allocations,omitempty"` // For frontend
}

type GoalAssetAllocation struct {
	ID        uuid.UUID `json:"id" db:"id"`
	GoalID    uuid.UUID `json:"goal_id" db:"goal_id"`
	AssetID   uuid.UUID `json:"asset_id" db:"asset_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LinkAssetToGoalRequest struct {
	AssetID uuid.UUID `json:"asset_id" validate:"required"`
}

type CreateFinancialGoalRequest struct {
	Name         string     `json:"name" validate:"required"`
	TargetAmount float64    `json:"target_amount" validate:"required,gt=0"`
	Priority     int        `json:"priority" validate:"required"`
	TargetDate   *time.Time `json:"target_date"`
}

type UpdateFinancialGoalRequest struct {
	Name         string     `json:"name"`
	TargetAmount float64    `json:"target_amount" validate:"omitempty,gt=0"`
	Priority     int        `json:"priority"`
	TargetDate   *time.Time `json:"target_date"`
	Status       string     `json:"status"`
}
