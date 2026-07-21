package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// InvestmentThesis represents the core reason for holding an asset.
// swagger:model InvestmentThesis
type InvestmentThesis struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	Ticker          string          `json:"ticker" db:"ticker"`
	CompanyName     string          `json:"company_name" db:"company_name"`
	Status          string          `json:"status" db:"status"` // active, archived, sold
	WhyIOwn         string          `json:"why_i_own" db:"why_i_own"`
	ThesisSummary   string          `json:"thesis_summary" db:"thesis_summary"`
	ThesisDetail    string          `json:"thesis_detail" db:"thesis_detail"`
	Moat            json.RawMessage `json:"moat" db:"moat" swaggertype:"array,string"` // Array of strings or complex objects
	Catalysts       json.RawMessage `json:"catalysts" db:"catalysts" swaggertype:"array,string"`
	Risks           json.RawMessage `json:"risks" db:"risks" swaggertype:"array,string"`
	KeyMetrics      json.RawMessage `json:"key_metrics" db:"key_metrics" swaggertype:"object"`
	SellConditions  json.RawMessage `json:"sell_conditions" db:"sell_conditions" swaggertype:"array,string"`
	ConvictionScore int             `json:"conviction_score" db:"conviction_score"` // 1-10
	QualityScore    int             `json:"quality_score" db:"quality_score"`       // 1-10
	ValuationScore  int             `json:"valuation_score" db:"valuation_score"`   // 1-10
	FairValue       float64         `json:"fair_value" db:"fair_value"`
	MarginOfSafety  float64         `json:"margin_of_safety" db:"margin_of_safety"`
	InitialDate     *time.Time      `json:"initial_date" db:"initial_date"`
	LastReviewed    *time.Time      `json:"last_reviewed" db:"last_reviewed"`
	Version         int             `json:"version" db:"version"`
	Notes           string          `json:"notes" db:"notes"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// ThesisGenerationRequest represents the payload to ask AI to generate a thesis
// swagger:model ThesisGenerationRequest
type ThesisGenerationRequest struct {
	Ticker       string `json:"ticker" validate:"required"`
	AssetType    string `json:"asset_type" validate:"required"` // stock, real_estate, crypto, etc.
	CompanyName  string `json:"company_name" validate:"required"`
	UserProvided string `json:"user_provided_context"`
}
