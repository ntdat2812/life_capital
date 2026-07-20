package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type AssetRepository interface {
	CreateAsset(ctx context.Context, asset *model.Asset) error
	GetAssetsByUserID(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedAssets, error)
	UpdateAsset(ctx context.Context, asset *model.Asset) error
	DeleteAsset(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type assetRepository struct {
	db *pgxpool.Pool
}

func NewAssetRepository(db *pgxpool.Pool) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) CreateAsset(ctx context.Context, asset *model.Asset) error {
	query := `
		INSERT INTO assets (
			user_id, category, name, ticker, quantity, avg_price, 
			current_price, current_value, cost_basis, notes, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id
	`
	now := time.Now()
	asset.CreatedAt = now
	asset.UpdatedAt = now
	asset.IsActive = true

	err := r.db.QueryRow(ctx, query,
		asset.UserID, asset.Category, asset.Name, asset.Ticker, asset.Quantity, asset.AvgPrice,
		asset.CurrentPrice, asset.CurrentValue, asset.CostBasis, asset.Notes, asset.IsActive, asset.CreatedAt, asset.UpdatedAt,
	).Scan(&asset.ID)

	return err
}

func (r *assetRepository) GetAssetsByUserID(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedAssets, error) {
	// First, count total items
	countQuery := `SELECT count(*) FROM assets WHERE user_id = $1 AND is_active = true`
	args := []interface{}{userID}
	if category != "" {
		countQuery += ` AND category = $2`
		args = append(args, category)
	}
	
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Then, fetch paginated items
	query := `
		SELECT id, user_id, category, name, ticker, quantity, avg_price, current_price, current_value, cost_basis, notes, is_active, created_at, updated_at
		FROM assets
		WHERE user_id = $1 AND is_active = true
	`
	args = []interface{}{userID}
	argIdx := 2
	if category != "" {
		query += ` AND category = $` + string(rune('0'+argIdx))
		args = append(args, category)
		argIdx++
	}

	orderClause := "current_value DESC"
	switch sort {
	case "value_asc":
		orderClause = "current_value ASC"
	case "name_asc":
		orderClause = "name ASC"
	case "name_desc":
		orderClause = "name DESC"
	}

	query += ` ORDER BY ` + orderClause
	if limit > 0 {
		query += ` LIMIT $` + string(rune('0'+argIdx)) + ` OFFSET $` + string(rune('0'+argIdx+1))
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []model.Asset = make([]model.Asset, 0)
	for rows.Next() {
		var a model.Asset
		err := rows.Scan(
			&a.ID, &a.UserID, &a.Category, &a.Name, &a.Ticker, &a.Quantity, &a.AvgPrice,
			&a.CurrentPrice, &a.CurrentValue, &a.CostBasis, &a.Notes, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}

	return &model.PaginatedAssets{
		Data:       assets,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *assetRepository) UpdateAsset(ctx context.Context, asset *model.Asset) error {
	query := `
		UPDATE assets SET 
			category = $1, name = $2, ticker = $3, quantity = $4, avg_price = $5, 
			current_price = $6, current_value = $7, cost_basis = $8, notes = $9, updated_at = $10
		WHERE id = $11 AND user_id = $12 AND is_active = true
	`
	asset.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		asset.Category, asset.Name, asset.Ticker, asset.Quantity, asset.AvgPrice,
		asset.CurrentPrice, asset.CurrentValue, asset.CostBasis, asset.Notes, asset.UpdatedAt,
		asset.ID, asset.UserID,
	)
	return err
}

func (r *assetRepository) DeleteAsset(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM assets WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}
