package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type MonthlyReviewRepository interface {
	Create(ctx context.Context, review *model.MonthlyReview) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.MonthlyReview, error)
	GetByMonth(ctx context.Context, userID uuid.UUID, reviewMonth string) (*model.MonthlyReview, error)
	GetLatestReview(ctx context.Context, userID uuid.UUID) (*model.MonthlyReview, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.MonthlyReview, error)
	Update(ctx context.Context, review *model.MonthlyReview) error
}

type monthlyReviewRepository struct {
	db *pgxpool.Pool
}

func NewMonthlyReviewRepository(db *pgxpool.Pool) MonthlyReviewRepository {
	return &monthlyReviewRepository{db: db}
}

func (r *monthlyReviewRepository) Create(ctx context.Context, review *model.MonthlyReview) error {
	query := `
		INSERT INTO monthly_reviews (user_id, review_month, status, new_investment_amount, portfolio_snapshot, net_worth_at_review, ai_health_assessment, ai_recommendations, ai_overall_summary, user_note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	snapshotJSON, _ := json.Marshal(review.PortfolioSnapshot)
	recommendationsJSON, _ := json.Marshal(review.AIRecommendations)

	return r.db.QueryRow(ctx, query,
		review.UserID,
		review.ReviewMonth,
		review.Status,
		review.NewInvestmentAmount,
		string(snapshotJSON),
		review.NetWorthAtReview,
		review.AIHealthAssessment,
		string(recommendationsJSON),
		review.AIOverallSummary,
		review.UserNote,
	).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
}

func (r *monthlyReviewRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.MonthlyReview, error) {
	query := `
		SELECT id, user_id, review_month, status, new_investment_amount, portfolio_snapshot, net_worth_at_review, ai_health_assessment, ai_recommendations, ai_overall_summary, user_note, created_at, updated_at
		FROM monthly_reviews
		WHERE id = $1 AND user_id = $2
	`
	var review model.MonthlyReview
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&review.ID,
		&review.UserID,
		&review.ReviewMonth,
		&review.Status,
		&review.NewInvestmentAmount,
		&review.PortfolioSnapshot,
		&review.NetWorthAtReview,
		&review.AIHealthAssessment,
		&review.AIRecommendations,
		&review.AIOverallSummary,
		&review.UserNote,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	return &review, err
}

func (r *monthlyReviewRepository) GetByMonth(ctx context.Context, userID uuid.UUID, reviewMonth string) (*model.MonthlyReview, error) {
	query := `
		SELECT id, user_id, review_month, status, new_investment_amount, portfolio_snapshot, net_worth_at_review, ai_health_assessment, ai_recommendations, ai_overall_summary, user_note, created_at, updated_at
		FROM monthly_reviews
		WHERE user_id = $1 AND review_month = $2
	`
	var review model.MonthlyReview
	err := r.db.QueryRow(ctx, query, userID, reviewMonth).Scan(
		&review.ID,
		&review.UserID,
		&review.ReviewMonth,
		&review.Status,
		&review.NewInvestmentAmount,
		&review.PortfolioSnapshot,
		&review.NetWorthAtReview,
		&review.AIHealthAssessment,
		&review.AIRecommendations,
		&review.AIOverallSummary,
		&review.UserNote,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	return &review, err
}

func (r *monthlyReviewRepository) GetLatestReview(ctx context.Context, userID uuid.UUID) (*model.MonthlyReview, error) {
	query := `
		SELECT id, user_id, review_month, status, new_investment_amount, portfolio_snapshot, net_worth_at_review, ai_health_assessment, ai_recommendations, ai_overall_summary, user_note, created_at, updated_at
		FROM monthly_reviews
		WHERE user_id = $1 AND status = 'completed'
		ORDER BY review_month DESC
		LIMIT 1
	`
	var review model.MonthlyReview
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&review.ID,
		&review.UserID,
		&review.ReviewMonth,
		&review.Status,
		&review.NewInvestmentAmount,
		&review.PortfolioSnapshot,
		&review.NetWorthAtReview,
		&review.AIHealthAssessment,
		&review.AIRecommendations,
		&review.AIOverallSummary,
		&review.UserNote,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil if no completed reviews exist
		}
		return nil, err
	}
	return &review, nil
}

func (r *monthlyReviewRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.MonthlyReview, error) {
	query := `
		SELECT id, user_id, review_month, status, new_investment_amount, portfolio_snapshot, net_worth_at_review, ai_health_assessment, ai_recommendations, ai_overall_summary, user_note, created_at, updated_at
		FROM monthly_reviews
		WHERE user_id = $1
		ORDER BY review_month DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.MonthlyReview
	for rows.Next() {
		var review model.MonthlyReview
		if err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ReviewMonth,
			&review.Status,
			&review.NewInvestmentAmount,
			&review.PortfolioSnapshot,
			&review.NetWorthAtReview,
			&review.AIHealthAssessment,
			&review.AIRecommendations,
			&review.AIOverallSummary,
			&review.UserNote,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

func (r *monthlyReviewRepository) Update(ctx context.Context, review *model.MonthlyReview) error {
	query := `
		UPDATE monthly_reviews
		SET status = $1, new_investment_amount = $2, portfolio_snapshot = $3, net_worth_at_review = $4, ai_health_assessment = $5, ai_recommendations = $6, ai_overall_summary = $7, user_note = $8, updated_at = NOW()
		WHERE id = $9 AND user_id = $10
		RETURNING updated_at
	`
	snapshotJSON, _ := json.Marshal(review.PortfolioSnapshot)
	recommendationsJSON, _ := json.Marshal(review.AIRecommendations)

	return r.db.QueryRow(ctx, query,
		review.Status,
		review.NewInvestmentAmount,
		string(snapshotJSON),
		review.NetWorthAtReview,
		review.AIHealthAssessment,
		string(recommendationsJSON),
		review.AIOverallSummary,
		review.UserNote,
		review.ID,
		review.UserID,
	).Scan(&review.UpdatedAt)
}
