package repository

import (
	"context"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvestorProfileRepository struct {
	db *pgxpool.Pool
}

func NewInvestorProfileRepository(db *pgxpool.Pool) *InvestorProfileRepository {
	return &InvestorProfileRepository{db: db}
}

func (r *InvestorProfileRepository) CreateProfile(ctx context.Context, profile *model.InvestorProfile) error {
	query := `
		INSERT INTO investor_profiles (
			user_id, version, status, date_of_birth, marital_status, risk_tolerance,
			risk_score, total_monthly_income, total_monthly_expense, fi_target_amount, life_constraints
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		profile.UserID, profile.Version, profile.Status, profile.DateOfBirth, profile.MaritalStatus,
		profile.RiskTolerance, profile.RiskScore, profile.TotalMonthlyIncome, profile.TotalMonthlyExpense,
		profile.FITargetAmount, profile.LifeConstraints,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

	return err
}

func (r *InvestorProfileRepository) GetActiveProfileByUserID(ctx context.Context, userID string) (*model.InvestorProfile, error) {
	query := `
		SELECT id, user_id, version, status, date_of_birth, marital_status, risk_tolerance,
			risk_score, total_monthly_income, total_monthly_expense, fi_target_amount, life_constraints,
			created_at, updated_at
		FROM investor_profiles
		WHERE user_id = $1 AND status = 'active'
		ORDER BY version DESC LIMIT 1
	`
	profile := &model.InvestorProfile{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Version, &profile.Status, &profile.DateOfBirth,
		&profile.MaritalStatus, &profile.RiskTolerance, &profile.RiskScore, &profile.TotalMonthlyIncome,
		&profile.TotalMonthlyExpense, &profile.FITargetAmount, &profile.LifeConstraints,
		&profile.CreatedAt, &profile.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return profile, nil
}

func (r *InvestorProfileRepository) SupersedePreviousProfiles(ctx context.Context, userID string) error {
	query := `
		UPDATE investor_profiles
		SET status = 'superseded', updated_at = NOW()
		WHERE user_id = $1 AND status = 'active'
	`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
