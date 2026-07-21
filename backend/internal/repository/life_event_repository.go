package repository

import (
	"context"
	"time"

	"github.com/datnguyen/life_capital/backend/internal/model"
)

type LifeEventRepository interface {
	Create(ctx context.Context, event *model.LifeEvent) error
	GetByUserID(ctx context.Context, userID string) ([]model.LifeEvent, error)
}

type lifeEventRepository struct {
	db DBTX
}

func NewLifeEventRepository(db DBTX) LifeEventRepository {
	return &lifeEventRepository{db: db}
}

func (r *lifeEventRepository) Create(ctx context.Context, event *model.LifeEvent) error {
	query := `
		INSERT INTO life_events (
			user_id, event_date, category, title, income_impact, expense_impact,
			ai_impact_analysis, triggered_profile_version, triggered_ips_version, requires_ips_update
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id, created_at
	`
	
	if event.EventDate.IsZero() {
		event.EventDate = time.Now()
	}

	err := r.db.QueryRow(ctx, query,
		event.UserID,
		event.EventDate,
		event.Category,
		event.Title,
		event.IncomeImpact,
		event.ExpenseImpact,
		event.AIImpactAnalysis,
		event.TriggeredProfileVersion,
		event.TriggeredIPSVersion,
		event.RequiresIPSUpdate,
	).Scan(&event.ID, &event.CreatedAt)

	return err
}

func (r *lifeEventRepository) GetByUserID(ctx context.Context, userID string) ([]model.LifeEvent, error) {
	query := `
		SELECT id, user_id, event_date, category, title, income_impact, expense_impact,
		       ai_impact_analysis, triggered_profile_version, triggered_ips_version, requires_ips_update, created_at
		FROM life_events
		WHERE user_id = $1
		ORDER BY event_date DESC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.LifeEvent
	for rows.Next() {
		var e model.LifeEvent
		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.EventDate,
			&e.Category,
			&e.Title,
			&e.IncomeImpact,
			&e.ExpenseImpact,
			&e.AIImpactAnalysis,
			&e.TriggeredProfileVersion,
			&e.TriggeredIPSVersion,
			&e.RequiresIPSUpdate,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
