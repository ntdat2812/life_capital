package model

import (
	"time"

	"github.com/google/uuid"
)

type AssetCategory string

const (
	AssetCategoryCash       AssetCategory = "cash"
	AssetCategoryDeposit    AssetCategory = "deposit"
	AssetCategoryGold       AssetCategory = "gold"
	AssetCategoryStock      AssetCategory = "stock"
	AssetCategoryFund       AssetCategory = "fund"
	AssetCategoryCrypto     AssetCategory = "crypto"
	AssetCategoryRealEstate AssetCategory = "real_estate"
)

type LiabilityCategory string

const (
	LiabilityCategoryMortgage     LiabilityCategory = "mortgage"
	LiabilityCategoryAutoLoan     LiabilityCategory = "auto_loan"
	LiabilityCategoryStudentLoan  LiabilityCategory = "student_loan"
	LiabilityCategoryCreditCard   LiabilityCategory = "credit_card"
	LiabilityCategoryPersonalLoan LiabilityCategory = "personal_loan"
	LiabilityCategoryOther        LiabilityCategory = "other"
)

type Asset struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	UserID       uuid.UUID     `json:"user_id" db:"user_id"`
	Category     AssetCategory `json:"category" db:"category"`
	Name         string        `json:"name" db:"name"`
	Ticker       *string       `json:"ticker,omitempty" db:"ticker"`
	Quantity     *float64      `json:"quantity,omitempty" db:"quantity"`
	AvgPrice     *float64      `json:"avg_price,omitempty" db:"avg_price"`
	CurrentPrice *float64      `json:"current_price,omitempty" db:"current_price"`
	CurrentValue float64       `json:"current_value" db:"current_value"`
	CostBasis    *float64      `json:"cost_basis,omitempty" db:"cost_basis"`
	Notes        *string       `json:"notes,omitempty" db:"notes"`
	IsActive     bool          `json:"is_active" db:"is_active"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}

type Liability struct {
	ID               uuid.UUID         `json:"id" db:"id"`
	UserID           uuid.UUID         `json:"user_id" db:"user_id"`
	Category         LiabilityCategory `json:"category" db:"category"`
	Name             string            `json:"name" db:"name"`
	RemainingBalance float64           `json:"remaining_balance" db:"remaining_balance"`
	InterestRate     *float64          `json:"interest_rate,omitempty" db:"interest_rate"`
	MonthlyPayment   float64           `json:"monthly_payment" db:"monthly_payment"`
	Lender           *string           `json:"lender,omitempty" db:"lender"`
	Notes            *string           `json:"notes,omitempty" db:"notes"`
	IsActive         bool              `json:"is_active" db:"is_active"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
}

type NetWorthSummary struct {
	TotalAssets      float64 `json:"total_assets"`
	TotalLiabilities float64 `json:"total_liabilities"`
	NetWorth         float64 `json:"net_worth"`
	BaseCurrency     string  `json:"base_currency"`
}

// DTOs
type CreateAssetRequest struct {
	Category     AssetCategory `json:"category" validate:"required"`
	Name         string        `json:"name" validate:"required"`
	Ticker       *string       `json:"ticker,omitempty"`
	Quantity     *float64      `json:"quantity,omitempty"`
	AvgPrice     *float64      `json:"avg_price,omitempty"`
	CurrentPrice *float64      `json:"current_price,omitempty"`
	CurrentValue float64       `json:"current_value" validate:"required"`
	CostBasis    *float64      `json:"cost_basis,omitempty"`
	Notes        *string       `json:"notes,omitempty"`
}

type CreateLiabilityRequest struct {
	Category         LiabilityCategory `json:"category" validate:"required"`
	Name             string            `json:"name" validate:"required"`
	RemainingBalance float64           `json:"remaining_balance" validate:"required"`
	InterestRate     *float64          `json:"interest_rate,omitempty"`
	MonthlyPayment   float64           `json:"monthly_payment"`
	Lender           *string           `json:"lender,omitempty"`
	Notes            *string           `json:"notes,omitempty"`
}

type UpdateAssetRequest struct {
	Category     AssetCategory `json:"category" validate:"required"`
	Name         string        `json:"name" validate:"required"`
	Ticker       *string       `json:"ticker,omitempty"`
	Quantity     *float64      `json:"quantity,omitempty"`
	AvgPrice     *float64      `json:"avg_price,omitempty"`
	CurrentPrice *float64      `json:"current_price,omitempty"`
	CurrentValue float64       `json:"current_value" validate:"required"`
	CostBasis    *float64      `json:"cost_basis,omitempty"`
	Notes        *string       `json:"notes,omitempty"`
}

type UpdateLiabilityRequest struct {
	Category         LiabilityCategory `json:"category" validate:"required"`
	Name             string            `json:"name" validate:"required"`
	RemainingBalance float64           `json:"remaining_balance" validate:"required"`
	InterestRate     *float64          `json:"interest_rate,omitempty"`
	MonthlyPayment   float64           `json:"monthly_payment"`
	Lender           *string           `json:"lender,omitempty"`
	Notes            *string           `json:"notes,omitempty"`
}

type PaginatedAssets struct {
	Data       []Asset `json:"data"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}

type PaginatedLiabilities struct {
	Data       []Liability `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
