package model

import (
	"time"

	"github.com/google/uuid"
)

// MonthlyReview represents a monthly snapshot and AI recommendation.
// swagger:model MonthlyReview
type MonthlyReview struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	UserID              uuid.UUID `json:"user_id" db:"user_id"`
	ReviewMonth         time.Time `json:"review_month" db:"review_month"` // e.g. "2023-10-01T00:00:00Z" from DB, we'll format it
	Status              string    `json:"status" db:"status"` // draft, completed
	NewInvestmentAmount float64   `json:"new_investment_amount" db:"new_investment_amount"`
	
	// Snapshots (JSONB)
	PortfolioSnapshot interface{} `json:"portfolio_snapshot" db:"portfolio_snapshot"`
	NetWorthAtReview  float64     `json:"net_worth_at_review" db:"net_worth_at_review"`
	
	// AI Generated content
	AIRecommendations interface{} `json:"ai_recommendations" db:"ai_recommendations"`
	AIOverallSummary  string      `json:"ai_overall_summary" db:"ai_overall_summary"`
	
	UserNote          *string     `json:"user_note" db:"user_note"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// MonthlyReviewRecommendation represents an individual action recommended by AI
type MonthlyReviewRecommendation struct {
	Ticker string  `json:"ticker"`
	Action string  `json:"action"` // "buy", "sell", "hold", "add", "reduce", "pay_debt"
	Amount float64 `json:"amount"` // The amount to invest or divest
	Reason string  `json:"reason"`
}

// MonthlyReviewRecommendationResponse represents the AI response structure
type MonthlyReviewRecommendationResponse struct {
	AIOverallSummary  string                        `json:"ai_overall_summary"`
	AIRecommendations []MonthlyReviewRecommendation `json:"ai_recommendations"`
}
