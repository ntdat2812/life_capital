package model

import (
	"time"

	"github.com/google/uuid"
)

type InvestorProfile struct {
	ID                          uuid.UUID              `json:"id" db:"id"`
	UserID                      uuid.UUID              `json:"user_id" db:"user_id"`
	Version                     int                    `json:"version" db:"version"`
	Status                      string                 `json:"status" db:"status"` // active, superseded
	DateOfBirth                 *time.Time             `json:"date_of_birth" db:"date_of_birth"`
	MaritalStatus               string                 `json:"marital_status" db:"marital_status"`
	RiskTolerance               string                 `json:"risk_tolerance" db:"risk_tolerance"`
	RiskScore                   int                    `json:"risk_score" db:"risk_score"`
	EssentialMonthlyExpense     float64                `json:"essential_monthly_expense" db:"essential_monthly_expense"`
	DiscretionaryMonthlyExpense float64                `json:"discretionary_monthly_expense" db:"discretionary_monthly_expense"`
	FITargetAmount              float64                `json:"fi_target_amount" db:"fi_target_amount"`
	CreatedAt                   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time              `json:"updated_at" db:"updated_at"`
}

type OnboardingRequest struct {
	ChatHistory []ChatMessage `json:"chat_history"`
}

type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // The text content
}

type UpdateProfileRequest struct {
	DateOfBirth                 *time.Time `json:"date_of_birth"`
	MaritalStatus               string     `json:"marital_status"`
	RiskTolerance               string     `json:"risk_tolerance"`
	RiskScore                   int        `json:"risk_score"`
	EssentialMonthlyExpense     float64    `json:"essential_monthly_expense"`
	DiscretionaryMonthlyExpense float64    `json:"discretionary_monthly_expense"`
	FITargetAmount              float64    `json:"fi_target_amount"`
}
