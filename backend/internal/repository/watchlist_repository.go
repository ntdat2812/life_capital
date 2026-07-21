package repository

import (
	"context"
	"fmt"
	"github.com/datnguyen/life_capital/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchlistRepository interface {
	Create(ctx context.Context, item *model.WatchlistItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.WatchlistItem, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.WatchlistItem, error)
	Update(ctx context.Context, item *model.WatchlistItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type watchlistRepository struct {
	db *pgxpool.Pool
}

func NewWatchlistRepository(db *pgxpool.Pool) WatchlistRepository {
	return &watchlistRepository{db: db}
}

func (r *watchlistRepository) Create(ctx context.Context, item *model.WatchlistItem) error {
	query := `
		INSERT INTO watchlist (
			id, user_id, ticker, company_name, thesis_id, added_date, 
			target_price, current_price, fair_value, quality_score, 
			status, priority, notes, ai_alert, last_ai_check
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING created_at, updated_at
	`
	
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	err := r.db.QueryRow(ctx, query,
		item.ID, item.UserID, item.Ticker, item.CompanyName, item.ThesisID,
		item.AddedDate, item.TargetPrice, item.CurrentPrice, item.FairValue,
		item.QualityScore, item.Status, item.Priority, item.Notes,
		item.AIAlert, item.LastAICheck,
	).Scan(&item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create watchlist item: %w", err)
	}
	return nil
}

func (r *watchlistRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.WatchlistItem, error) {
	query := `
		SELECT 
			id, user_id, ticker, company_name, thesis_id, added_date, 
			target_price, current_price, fair_value, quality_score, 
			status, priority, notes, ai_alert, last_ai_check, 
			created_at, updated_at
		FROM watchlist
		WHERE id = $1
	`
	
	var w model.WatchlistItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&w.ID, &w.UserID, &w.Ticker, &w.CompanyName, &w.ThesisID,
		&w.AddedDate, &w.TargetPrice, &w.CurrentPrice, &w.FairValue,
		&w.QualityScore, &w.Status, &w.Priority, &w.Notes,
		&w.AIAlert, &w.LastAICheck, &w.CreatedAt, &w.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist item: %w", err)
	}
	return &w, nil
}

func (r *watchlistRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.WatchlistItem, error) {
	query := `
		SELECT 
			id, user_id, ticker, company_name, thesis_id, added_date, 
			target_price, current_price, fair_value, quality_score, 
			status, priority, notes, ai_alert, last_ai_check, 
			created_at, updated_at
		FROM watchlist
		WHERE user_id = $1
		ORDER BY priority DESC, created_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query watchlist: %w", err)
	}
	defer rows.Close()

	var items []*model.WatchlistItem
	for rows.Next() {
		var w model.WatchlistItem
		err := rows.Scan(
			&w.ID, &w.UserID, &w.Ticker, &w.CompanyName, &w.ThesisID,
			&w.AddedDate, &w.TargetPrice, &w.CurrentPrice, &w.FairValue,
			&w.QualityScore, &w.Status, &w.Priority, &w.Notes,
			&w.AIAlert, &w.LastAICheck, &w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan watchlist row: %w", err)
		}
		items = append(items, &w)
	}

	return items, nil
}

func (r *watchlistRepository) Update(ctx context.Context, item *model.WatchlistItem) error {
	query := `
		UPDATE watchlist SET
			company_name = $1, thesis_id = $2, target_price = $3, 
			current_price = $4, fair_value = $5, quality_score = $6, 
			status = $7, priority = $8, notes = $9, ai_alert = $10, 
			last_ai_check = $11, updated_at = NOW()
		WHERE id = $12 AND user_id = $13
		RETURNING updated_at
	`
	
	err := r.db.QueryRow(ctx, query,
		item.CompanyName, item.ThesisID, item.TargetPrice,
		item.CurrentPrice, item.FairValue, item.QualityScore,
		item.Status, item.Priority, item.Notes, item.AIAlert,
		item.LastAICheck, item.ID, item.UserID,
	).Scan(&item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update watchlist item: %w", err)
	}
	return nil
}

func (r *watchlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM watchlist WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete watchlist item: %w", err)
	}
	return nil
}
