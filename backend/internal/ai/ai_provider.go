package ai

import "context"

type ExtractionResult struct {
	DateOfBirth         *string                `json:"date_of_birth"`
	RiskScore           int                    `json:"risk_score"`
	RiskTolerance       string                 `json:"risk_tolerance"`
	MaritalStatus       string                 `json:"marital_status"`
	TotalMonthlyIncome          float64                `json:"total_monthly_income"`
	EssentialMonthlyExpense     float64                `json:"essential_monthly_expense"`
	DiscretionaryMonthlyExpense float64                `json:"discretionary_monthly_expense"`
	FITargetAmount              float64                `json:"fi_target_amount"`
}

type AIProvider interface {
	ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error)
}
