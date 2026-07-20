package ai

import (
	"context"
	"fmt"
)

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

// ExecuteWithFallback is a generic helper to run any AIProvider method with a fallback mechanism.
func ExecuteWithFallback[T any](providers []AIProvider, action func(AIProvider) (T, error)) (T, error) {
	var zero T
	if len(providers) == 0 {
		return zero, fmt.Errorf("no AI providers configured")
	}

	var lastErr error
	for i, provider := range providers {
		result, err := action(provider)
		if err == nil {
			return result, nil
		}
		fmt.Printf("Warning: AI provider (index %d) failed, trying next: %v\n", i, err)
		lastErr = err
	}

	return zero, fmt.Errorf("all AI providers failed, last error: %v", lastErr)
}
