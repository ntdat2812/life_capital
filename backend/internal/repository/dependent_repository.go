package repository

import (
	"context"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DependentRepository struct {
	db *pgxpool.Pool
}

func NewDependentRepository(db *pgxpool.Pool) *DependentRepository {
	return &DependentRepository{db: db}
}

func (r *DependentRepository) CreateDependent(ctx context.Context, dep *model.Dependent) error {
	query := `
		INSERT INTO dependents (
			user_id, name, relationship, date_of_birth, monthly_cost, notes
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, is_active, added_date, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		dep.UserID, dep.Name, dep.Relationship, dep.DateOfBirth,
		dep.MonthlyCost, dep.Notes,
	).Scan(&dep.ID, &dep.IsActive, &dep.AddedDate, &dep.CreatedAt, &dep.UpdatedAt)
}

func (r *DependentRepository) GetDependentsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Dependent, error) {
	query := `
		SELECT id, user_id, name, relationship, date_of_birth, is_active, monthly_cost, notes, added_date, created_at, updated_at
		FROM dependents
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []*model.Dependent
	for rows.Next() {
		dep := &model.Dependent{}
		err := rows.Scan(
			&dep.ID, &dep.UserID, &dep.Name, &dep.Relationship, &dep.DateOfBirth,
			&dep.IsActive, &dep.MonthlyCost, &dep.Notes, &dep.AddedDate,
			&dep.CreatedAt, &dep.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}
	return deps, nil
}

func (r *DependentRepository) UpdateDependent(ctx context.Context, dep *model.Dependent) error {
	query := `
		UPDATE dependents
		SET name = $1, relationship = $2, date_of_birth = $3, is_active = $4,
		    monthly_cost = $5, notes = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
	`
	cmdTag, err := r.db.Exec(ctx, query,
		dep.Name, dep.Relationship, dep.DateOfBirth, dep.IsActive,
		dep.MonthlyCost, dep.Notes, dep.ID, dep.UserID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *DependentRepository) DeleteDependent(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM dependents WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
