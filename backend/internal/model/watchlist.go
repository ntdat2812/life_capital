package model

import (
	"time"

	"github.com/google/uuid"
)

// WatchlistItem represents an asset the user is tracking.
// swagger:model WatchlistItem
type WatchlistItem struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Ticker       string     `json:"ticker" db:"ticker"`
	CompanyName  string     `json:"company_name" db:"company_name"`
	ThesisID     *uuid.UUID `json:"thesis_id" db:"thesis_id"` // Nullable
	AddedDate    *time.Time `json:"added_date" db:"added_date"`
	TargetPrice  float64    `json:"target_price" db:"target_price"`
	CurrentPrice float64    `json:"current_price" db:"current_price"`
	FairValue    float64    `json:"fair_value" db:"fair_value"`
	QualityScore int        `json:"quality_score" db:"quality_score"`
	Status       string     `json:"status" db:"status"` // watching, bought, dropped
	Priority     int        `json:"priority" db:"priority"` // 1-10 (higher is higher priority)
	Notes        string     `json:"notes" db:"notes"`
	AIAlert      string     `json:"ai_alert" db:"ai_alert"`
	LastAICheck  *time.Time `json:"last_ai_check" db:"last_ai_check"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
