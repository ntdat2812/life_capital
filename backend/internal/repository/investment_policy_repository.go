package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type InvestmentPolicyRepository interface {
	Create(ctx context.Context, policy *model.InvestmentPolicy) error
	GetLatestByUserID(ctx context.Context, userID string) (*model.InvestmentPolicy, error)
	Update(ctx context.Context, policy *model.InvestmentPolicy) error
}

type investmentPolicyRepository struct {
	db *pgxpool.Pool
}

func NewInvestmentPolicyRepository(db *pgxpool.Pool) InvestmentPolicyRepository {
	return &investmentPolicyRepository{db: db}
}

func (r *investmentPolicyRepository) Create(ctx context.Context, p *model.InvestmentPolicy) error {
	query := `
		INSERT INTO investment_policies (
			user_id, version, profile_version, trigger_event_id, 
			target_allocation, detailed_strategy, is_ai_recommended, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		p.UserID, p.Version, p.ProfileVersion, p.TriggerEventID,
		p.TargetAllocation, p.DetailedStrategy, p.IsAIRecommended, p.Status,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *investmentPolicyRepository) GetLatestByUserID(ctx context.Context, userID string) (*model.InvestmentPolicy, error) {
	query := `
		SELECT id, user_id, version, profile_version, trigger_event_id, 
			target_allocation, detailed_strategy, is_ai_recommended, status, created_at, updated_at
		FROM investment_policies
		WHERE user_id = $1
		ORDER BY version DESC
		LIMIT 1
	`
	var p model.InvestmentPolicy
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.Version, &p.ProfileVersion, &p.TriggerEventID,
		&p.TargetAllocation, &p.DetailedStrategy, &p.IsAIRecommended, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *investmentPolicyRepository) Update(ctx context.Context, p *model.InvestmentPolicy) error {
	query := `
		UPDATE investment_policies
		SET target_allocation = $1, detailed_strategy = $2, 
		    is_ai_recommended = $3, status = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at
	`
	return r.db.QueryRow(ctx, query,
		p.TargetAllocation, p.DetailedStrategy, 
		p.IsAIRecommended, p.Status, 
		p.ID, p.UserID,
	).Scan(&p.UpdatedAt)
}
