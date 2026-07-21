package repository

import (
	"context"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/jackc/pgx/v5"
)

type InvestorProfileRepository struct {
	db DBTX
}

func NewInvestorProfileRepository(db DBTX) *InvestorProfileRepository {
	return &InvestorProfileRepository{db: db}
}

func (r *InvestorProfileRepository) CreateProfile(ctx context.Context, profile *model.InvestorProfile) error {
	query := `
		INSERT INTO investor_profiles (
			user_id, version, status, date_of_birth, marital_status, risk_tolerance,
			risk_score, essential_monthly_expense, discretionary_monthly_expense, fi_target_amount
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		profile.UserID, profile.Version, profile.Status, profile.DateOfBirth, profile.MaritalStatus,
		profile.RiskTolerance, profile.RiskScore, profile.EssentialMonthlyExpense, profile.DiscretionaryMonthlyExpense,
		profile.FITargetAmount,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

	return err
}

func (r *InvestorProfileRepository) GetActiveProfileByUserID(ctx context.Context, userID string) (*model.InvestorProfile, error) {
	query := `
		SELECT id, user_id, version, status, date_of_birth, marital_status, risk_tolerance,
			risk_score, essential_monthly_expense, discretionary_monthly_expense, fi_target_amount,
			created_at, updated_at
		FROM investor_profiles
		WHERE user_id = $1 AND status = 'active'
		ORDER BY version DESC LIMIT 1
	`
	profile := &model.InvestorProfile{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Version, &profile.Status, &profile.DateOfBirth,
		&profile.MaritalStatus, &profile.RiskTolerance, &profile.RiskScore, &profile.EssentialMonthlyExpense,
		&profile.DiscretionaryMonthlyExpense, &profile.FITargetAmount,
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

func (r *InvestorProfileRepository) UpdateProfile(ctx context.Context, profile *model.InvestorProfile) error {
	query := `
		UPDATE investor_profiles
		SET date_of_birth = $1, marital_status = $2, risk_tolerance = $3, risk_score = $4,
		    essential_monthly_expense = $5, discretionary_monthly_expense = $6, fi_target_amount = $7, updated_at = NOW()
		WHERE id = $8 AND user_id = $9 AND status = 'active'
	`
	_, err := r.db.Exec(ctx, query,
		profile.DateOfBirth, profile.MaritalStatus, profile.RiskTolerance, profile.RiskScore,
		profile.EssentialMonthlyExpense, profile.DiscretionaryMonthlyExpense, profile.FITargetAmount,
		profile.ID, profile.UserID,
	)
	return err
}
