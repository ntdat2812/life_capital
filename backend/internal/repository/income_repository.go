package repository

import (
	"context"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncomeRepository struct {
	db *pgxpool.Pool
}

func NewIncomeRepository(db *pgxpool.Pool) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) CreateIncomeStream(ctx context.Context, income *model.IncomeStream) error {
	query := `
		INSERT INTO income_streams (
			user_id, name, type, is_passive, amount, frequency, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, is_active, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		income.UserID, income.Name, income.Type, income.IsPassive,
		income.Amount, income.Frequency, income.Notes,
	).Scan(&income.ID, &income.IsActive, &income.CreatedAt, &income.UpdatedAt)
}

func (r *IncomeRepository) GetIncomeStreamsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.IncomeStream, error) {
	query := `
		SELECT id, user_id, name, type, is_passive, amount, frequency, is_active, start_date, end_date, notes, created_at, updated_at
		FROM income_streams
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streams []*model.IncomeStream
	for rows.Next() {
		stream := &model.IncomeStream{}
		err := rows.Scan(
			&stream.ID, &stream.UserID, &stream.Name, &stream.Type, &stream.IsPassive,
			&stream.Amount, &stream.Frequency, &stream.IsActive, &stream.StartDate,
			&stream.EndDate, &stream.Notes, &stream.CreatedAt, &stream.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

func (r *IncomeRepository) UpdateIncomeStream(ctx context.Context, income *model.IncomeStream) error {
	query := `
		UPDATE income_streams
		SET name = $1, type = $2, is_passive = $3, amount = $4, frequency = $5,
		    is_active = $6, notes = $7, updated_at = NOW()
		WHERE id = $8 AND user_id = $9
	`
	cmdTag, err := r.db.Exec(ctx, query,
		income.Name, income.Type, income.IsPassive, income.Amount, income.Frequency,
		income.IsActive, income.Notes, income.ID, income.UserID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *IncomeRepository) DeleteIncomeStream(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM income_streams WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
