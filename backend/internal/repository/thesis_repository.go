package repository

import (
	"context"
	"fmt"
	"github.com/datnguyen/life_capital/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThesisRepository interface {
	Create(ctx context.Context, thesis *model.InvestmentThesis) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.InvestmentThesis, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.InvestmentThesis, error)
	GetByTicker(ctx context.Context, userID uuid.UUID, ticker string) (*model.InvestmentThesis, error)
	Update(ctx context.Context, thesis *model.InvestmentThesis) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type thesisRepository struct {
	db *pgxpool.Pool
}

func NewThesisRepository(db *pgxpool.Pool) ThesisRepository {
	return &thesisRepository{db: db}
}

func (r *thesisRepository) Create(ctx context.Context, thesis *model.InvestmentThesis) error {
	query := `
		INSERT INTO investment_theses (
			id, user_id, ticker, company_name, status, why_i_own, 
			thesis_summary, thesis_detail, moat, catalysts, risks, 
			key_metrics, sell_conditions, conviction_score, quality_score, 
			valuation_score, fair_value, margin_of_safety, initial_date, 
			last_reviewed, version, notes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 
			$14, $15, $16, $17, $18, $19, $20, $21, $22
		) RETURNING created_at, updated_at
	`
	
	if thesis.ID == uuid.Nil {
		thesis.ID = uuid.New()
	}

	err := r.db.QueryRow(ctx, query,
		thesis.ID, thesis.UserID, thesis.Ticker, thesis.CompanyName, thesis.Status,
		thesis.WhyIOwn, thesis.ThesisSummary, thesis.ThesisDetail,
		thesis.Moat, thesis.Catalysts, thesis.Risks, thesis.KeyMetrics,
		thesis.SellConditions, thesis.ConvictionScore, thesis.QualityScore,
		thesis.ValuationScore, thesis.FairValue, thesis.MarginOfSafety,
		thesis.InitialDate, thesis.LastReviewed, thesis.Version, thesis.Notes,
	).Scan(&thesis.CreatedAt, &thesis.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create investment thesis: %w", err)
	}
	return nil
}

func (r *thesisRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.InvestmentThesis, error) {
	query := `
		SELECT 
			id, user_id, ticker, company_name, status, why_i_own, 
			thesis_summary, thesis_detail, moat, catalysts, risks, 
			key_metrics, sell_conditions, conviction_score, quality_score, 
			valuation_score, fair_value, margin_of_safety, initial_date, 
			last_reviewed, version, notes, created_at, updated_at
		FROM investment_theses
		WHERE id = $1
	`
	
	var t model.InvestmentThesis
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.Ticker, &t.CompanyName, &t.Status, &t.WhyIOwn,
		&t.ThesisSummary, &t.ThesisDetail, &t.Moat, &t.Catalysts, &t.Risks,
		&t.KeyMetrics, &t.SellConditions, &t.ConvictionScore, &t.QualityScore,
		&t.ValuationScore, &t.FairValue, &t.MarginOfSafety, &t.InitialDate,
		&t.LastReviewed, &t.Version, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get investment thesis: %w", err)
	}
	return &t, nil
}

func (r *thesisRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.InvestmentThesis, error) {
	query := `
		SELECT 
			id, user_id, ticker, company_name, status, why_i_own, 
			thesis_summary, thesis_detail, moat, catalysts, risks, 
			key_metrics, sell_conditions, conviction_score, quality_score, 
			valuation_score, fair_value, margin_of_safety, initial_date, 
			last_reviewed, version, notes, created_at, updated_at
		FROM investment_theses
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query investment theses: %w", err)
	}
	defer rows.Close()

	var theses []*model.InvestmentThesis
	for rows.Next() {
		var t model.InvestmentThesis
		err := rows.Scan(
			&t.ID, &t.UserID, &t.Ticker, &t.CompanyName, &t.Status, &t.WhyIOwn,
			&t.ThesisSummary, &t.ThesisDetail, &t.Moat, &t.Catalysts, &t.Risks,
			&t.KeyMetrics, &t.SellConditions, &t.ConvictionScore, &t.QualityScore,
			&t.ValuationScore, &t.FairValue, &t.MarginOfSafety, &t.InitialDate,
			&t.LastReviewed, &t.Version, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan thesis row: %w", err)
		}
		theses = append(theses, &t)
	}

	return theses, nil
}

func (r *thesisRepository) GetByTicker(ctx context.Context, userID uuid.UUID, ticker string) (*model.InvestmentThesis, error) {
	query := `
		SELECT 
			id, user_id, ticker, company_name, status, why_i_own, 
			thesis_summary, thesis_detail, moat, catalysts, risks, 
			key_metrics, sell_conditions, conviction_score, quality_score, 
			valuation_score, fair_value, margin_of_safety, initial_date, 
			last_reviewed, version, notes, created_at, updated_at
		FROM investment_theses
		WHERE user_id = $1 AND ticker = $2
		ORDER BY created_at DESC LIMIT 1
	`
	
	var t model.InvestmentThesis
	err := r.db.QueryRow(ctx, query, userID, ticker).Scan(
		&t.ID, &t.UserID, &t.Ticker, &t.CompanyName, &t.Status, &t.WhyIOwn,
		&t.ThesisSummary, &t.ThesisDetail, &t.Moat, &t.Catalysts, &t.Risks,
		&t.KeyMetrics, &t.SellConditions, &t.ConvictionScore, &t.QualityScore,
		&t.ValuationScore, &t.FairValue, &t.MarginOfSafety, &t.InitialDate,
		&t.LastReviewed, &t.Version, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get thesis by ticker: %w", err)
	}
	return &t, nil
}

func (r *thesisRepository) Update(ctx context.Context, thesis *model.InvestmentThesis) error {
	query := `
		UPDATE investment_theses SET
			company_name = $1, status = $2, why_i_own = $3, 
			thesis_summary = $4, thesis_detail = $5, moat = $6, catalysts = $7, risks = $8, 
			key_metrics = $9, sell_conditions = $10, conviction_score = $11, quality_score = $12, 
			valuation_score = $13, fair_value = $14, margin_of_safety = $15, 
			last_reviewed = $16, version = $17, notes = $18, updated_at = NOW()
		WHERE id = $19 AND user_id = $20
		RETURNING updated_at
	`
	
	err := r.db.QueryRow(ctx, query,
		thesis.CompanyName, thesis.Status, thesis.WhyIOwn, thesis.ThesisSummary,
		thesis.ThesisDetail, thesis.Moat, thesis.Catalysts, thesis.Risks,
		thesis.KeyMetrics, thesis.SellConditions, thesis.ConvictionScore,
		thesis.QualityScore, thesis.ValuationScore, thesis.FairValue,
		thesis.MarginOfSafety, thesis.LastReviewed, thesis.Version,
		thesis.Notes, thesis.ID, thesis.UserID,
	).Scan(&thesis.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update investment thesis: %w", err)
	}
	return nil
}

func (r *thesisRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM investment_theses WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete investment thesis: %w", err)
	}
	return nil
}
