package ai

import "context"

type ExtractionResult struct {
	DateOfBirth         *string                `json:"date_of_birth"`
	RiskScore           int                    `json:"risk_score"`
	RiskTolerance       string                 `json:"risk_tolerance"`
	MaritalStatus       string                 `json:"marital_status"`
	TotalMonthlyIncome  float64                `json:"total_monthly_income"`
	TotalMonthlyExpense float64                `json:"total_monthly_expense"`
	FITargetAmount      float64                `json:"fi_target_amount"`
	LifeConstraints     map[string]interface{} `json:"life_constraints"`
}

type AIProvider interface {
	ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error)
}
