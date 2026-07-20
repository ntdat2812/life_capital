package model

import (
	"time"

	"github.com/google/uuid"
)

type IncomeStream struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Name      string     `json:"name" db:"name"`
	Type      string     `json:"type" db:"type"`
	IsPassive bool       `json:"is_passive" db:"is_passive"`
	Amount    float64    `json:"amount" db:"amount"`
	Frequency string     `json:"frequency" db:"frequency"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	StartDate *time.Time `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date" db:"end_date"`
	Notes     string     `json:"notes" db:"notes"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateIncomeStreamRequest struct {
	Name      string     `json:"name" validate:"required"`
	Type      string     `json:"type" validate:"required"`
	IsPassive bool       `json:"is_passive"`
	Amount    float64    `json:"amount" validate:"required,gt=0"`
	Frequency string     `json:"frequency"` // e.g., "monthly", "yearly"
	Notes     string     `json:"notes"`
}

type UpdateIncomeStreamRequest struct {
	Name      string     `json:"name" validate:"required"`
	Type      string     `json:"type" validate:"required"`
	IsPassive bool       `json:"is_passive"`
	Amount    float64    `json:"amount" validate:"required,gt=0"`
	Frequency string     `json:"frequency"`
	IsActive  bool       `json:"is_active"`
	Notes     string     `json:"notes"`
}

type Dependent struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Name         string     `json:"name" db:"name"`
	Relationship string     `json:"relationship" db:"relationship"`
	DateOfBirth  *time.Time `json:"date_of_birth" db:"date_of_birth"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	MonthlyCost  float64    `json:"monthly_cost" db:"monthly_cost"`
	Notes        string     `json:"notes" db:"notes"`
	AddedDate    *time.Time `json:"added_date" db:"added_date"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateDependentRequest struct {
	Name         string     `json:"name" validate:"required"`
	Relationship string     `json:"relationship" validate:"required"`
	DateOfBirth  *time.Time `json:"date_of_birth"`
	MonthlyCost  float64    `json:"monthly_cost" validate:"gte=0"`
	Notes        string     `json:"notes"`
}

type UpdateDependentRequest struct {
	Name         string     `json:"name" validate:"required"`
	Relationship string     `json:"relationship" validate:"required"`
	DateOfBirth  *time.Time `json:"date_of_birth"`
	IsActive     bool       `json:"is_active"`
	MonthlyCost  float64    `json:"monthly_cost" validate:"gte=0"`
	Notes        string     `json:"notes"`
}
