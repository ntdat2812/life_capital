package model

import (
	"time"

	"github.com/google/uuid"
)

type LifeEvent struct {
	ID                      uuid.UUID `json:"id"`
	UserID                  uuid.UUID `json:"user_id"`
	EventDate               time.Time `json:"event_date"`
	Category                string    `json:"category"` // Uses life_event_category enum
	Title                   string    `json:"title"`
	IncomeImpact            float64   `json:"income_impact"`
	ExpenseImpact           float64   `json:"expense_impact"`
	AIImpactAnalysis        *string   `json:"ai_impact_analysis,omitempty"`
	TriggeredProfileVersion *int      `json:"triggered_profile_version,omitempty"`
	TriggeredIPSVersion     *int      `json:"triggered_ips_version,omitempty"`
	RequiresIPSUpdate       bool      `json:"requires_ips_update"`
	CreatedAt               time.Time `json:"created_at"`
}
