package model

import (
	"encoding/json"
	"time"
)

type InvestmentPolicy struct {
	ID               string          `json:"id" db:"id"`
	UserID           string          `json:"user_id" db:"user_id"`
	Version          int             `json:"version" db:"version"`
	ProfileVersion   int             `json:"profile_version" db:"profile_version"`
	TriggerEventID   *string         `json:"trigger_event_id" db:"trigger_event_id"`
	TargetAllocation json.RawMessage `json:"target_allocation" db:"target_allocation" swaggertype:"object"`
	DetailedStrategy string          `json:"detailed_strategy" db:"detailed_strategy"`
	IsAIRecommended  bool            `json:"is_ai_recommended" db:"is_ai_recommended"`
	Status           string          `json:"status" db:"status"` // 'draft', 'active', 'superseded'
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}
