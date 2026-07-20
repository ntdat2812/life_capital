package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type LiabilityRepository interface {
	CreateLiability(ctx context.Context, liability *model.Liability) error
	GetLiabilitiesByUserID(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedLiabilities, error)
	UpdateLiability(ctx context.Context, liability *model.Liability) error
	DeleteLiability(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type liabilityRepository struct {
	db *pgxpool.Pool
}

func NewLiabilityRepository(db *pgxpool.Pool) LiabilityRepository {
	return &liabilityRepository{db: db}
}

func (r *liabilityRepository) CreateLiability(ctx context.Context, liability *model.Liability) error {
	query := `
		INSERT INTO liabilities (
			user_id, category, name, remaining_balance, interest_rate, 
			monthly_payment, lender, notes, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id
	`
	now := time.Now()
	liability.CreatedAt = now
	liability.UpdatedAt = now
	liability.IsActive = true

	err := r.db.QueryRow(ctx, query,
		liability.UserID, liability.Category, liability.Name, liability.RemainingBalance,
		liability.InterestRate, liability.MonthlyPayment, liability.Lender, liability.Notes,
		liability.IsActive, liability.CreatedAt, liability.UpdatedAt,
	).Scan(&liability.ID)

	return err
}

func (r *liabilityRepository) GetLiabilitiesByUserID(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedLiabilities, error) {
	// Count total items
	countQuery := `SELECT count(*) FROM liabilities WHERE user_id = $1 AND is_active = true`
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

	// Fetch paginated items
	query := `
		SELECT id, user_id, category, name, remaining_balance, interest_rate, monthly_payment, lender, notes, is_active, created_at, updated_at
		FROM liabilities
		WHERE user_id = $1 AND is_active = true
	`
	args = []interface{}{userID}
	argIdx := 2
	if category != "" {
		query += ` AND category = $` + string(rune('0'+argIdx))
		args = append(args, category)
		argIdx++
	}

	orderClause := "remaining_balance DESC"
	switch sort {
	case "value_asc":
		orderClause = "remaining_balance ASC"
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

	var liabilities []model.Liability = make([]model.Liability, 0)
	for rows.Next() {
		var l model.Liability
		err := rows.Scan(
			&l.ID, &l.UserID, &l.Category, &l.Name, &l.RemainingBalance,
			&l.InterestRate, &l.MonthlyPayment, &l.Lender, &l.Notes,
			&l.IsActive, &l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		liabilities = append(liabilities, l)
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}

	return &model.PaginatedLiabilities{
		Data:       liabilities,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *liabilityRepository) UpdateLiability(ctx context.Context, liability *model.Liability) error {
	query := `
		UPDATE liabilities SET 
			category = $1, name = $2, remaining_balance = $3, interest_rate = $4, 
			monthly_payment = $5, lender = $6, notes = $7, updated_at = $8
		WHERE id = $9 AND user_id = $10 AND is_active = true
	`
	liability.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		liability.Category, liability.Name, liability.RemainingBalance, liability.InterestRate,
		liability.MonthlyPayment, liability.Lender, liability.Notes, liability.UpdatedAt,
		liability.ID, liability.UserID,
	)
	return err
}

func (r *liabilityRepository) DeleteLiability(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM liabilities WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}
