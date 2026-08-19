package repository

import (
	"context"
	"fmt"

	"time"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/google/uuid"
)

type FinancialGoalRepository interface {
	Create(ctx context.Context, goal *model.FinancialGoal) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.FinancialGoal, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.FinancialGoal, error)
	Update(ctx context.Context, goal *model.FinancialGoal) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	LinkAsset(ctx context.Context, goalID uuid.UUID, assetID uuid.UUID) error
	UnlinkAsset(ctx context.Context, goalID uuid.UUID, assetID uuid.UUID) error
}

type financialGoalRepository struct {
	db DBTX
}

func NewFinancialGoalRepository(db DBTX) FinancialGoalRepository {
	return &financialGoalRepository{db: db}
}

func (r *financialGoalRepository) Create(ctx context.Context, goal *model.FinancialGoal) error {
	query := `
		INSERT INTO financial_goals (user_id, name, target_amount, priority, target_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		goal.UserID, goal.Name, goal.TargetAmount, goal.Priority, goal.TargetDate, goal.Status,
	).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create financial goal: %w", err)
	}
	return nil
}

func (r *financialGoalRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.FinancialGoal, error) {
	query := `
		SELECT id, user_id, name, target_amount, priority, target_date, status, created_at, updated_at
		FROM financial_goals
		WHERE id = $1 AND user_id = $2
	`
	var goal model.FinancialGoal
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&goal.ID, &goal.UserID, &goal.Name, &goal.TargetAmount, &goal.Priority,
		&goal.TargetDate, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get financial goal: %w", err)
	}
	return &goal, nil
}

func (r *financialGoalRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.FinancialGoal, error) {
	query := `
		SELECT id, user_id, name, target_amount, priority, target_date, status, created_at, updated_at
		FROM financial_goals
		WHERE user_id = $1
		ORDER BY priority ASC, created_at ASC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list financial goals: %w", err)
	}
	defer rows.Close()

	var goals []model.FinancialGoal
	for rows.Next() {
		var goal model.FinancialGoal
		if err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Name, &goal.TargetAmount, &goal.Priority,
			&goal.TargetDate, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan financial goal: %w", err)
		}
		goals = append(goals, goal)
	}

	// Fetch allocations
	if len(goals) > 0 {
		allocQuery := `
			SELECT id, goal_id, asset_id, created_at
			FROM goal_asset_allocations
			WHERE goal_id = ANY($1::uuid[])
		`
		var goalIDs []string
		for _, g := range goals {
			goalIDs = append(goalIDs, g.ID.String())
		}

		allocRows, err := r.db.Query(ctx, allocQuery, goalIDs)
		if err == nil {
			defer allocRows.Close()
			allocMap := make(map[uuid.UUID][]model.GoalAssetAllocation)
			for allocRows.Next() {
				var alloc model.GoalAssetAllocation
				var createdAt *time.Time
				if err := allocRows.Scan(&alloc.ID, &alloc.GoalID, &alloc.AssetID, &createdAt); err == nil {
					if createdAt != nil {
						alloc.CreatedAt = *createdAt
					}
					allocMap[alloc.GoalID] = append(allocMap[alloc.GoalID], alloc)
				} else {
					fmt.Printf("Error scanning allocation: %v\n", err)
				}
			}
			for i := range goals {
				goals[i].Allocations = allocMap[goals[i].ID]
			}
		} else {
			fmt.Printf("Error querying allocations: %v\n", err)
		}
	}

	return goals, nil
}

func (r *financialGoalRepository) Update(ctx context.Context, goal *model.FinancialGoal) error {
	query := `
		UPDATE financial_goals
		SET name = $1, target_amount = $2, priority = $3, target_date = $4, status = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND user_id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		goal.Name, goal.TargetAmount, goal.Priority, goal.TargetDate, goal.Status, goal.ID, goal.UserID,
	).Scan(&goal.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update financial goal: %w", err)
	}
	return nil
}

func (r *financialGoalRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM financial_goals WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete financial goal: %w", err)
	}
	return nil
}

func (r *financialGoalRepository) LinkAsset(ctx context.Context, goalID uuid.UUID, assetID uuid.UUID) error {
	query := `
		INSERT INTO goal_asset_allocations (goal_id, asset_id)
		VALUES ($1, $2)
		ON CONFLICT (goal_id, asset_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, goalID, assetID)
	if err != nil {
		return fmt.Errorf("failed to link asset to goal: %w", err)
	}
	return nil
}

func (r *financialGoalRepository) UnlinkAsset(ctx context.Context, goalID uuid.UUID, assetID uuid.UUID) error {
	query := `DELETE FROM goal_asset_allocations WHERE goal_id = $1 AND asset_id = $2`
	_, err := r.db.Exec(ctx, query, goalID, assetID)
	if err != nil {
		return fmt.Errorf("failed to unlink asset from goal: %w", err)
	}
	return nil
}
