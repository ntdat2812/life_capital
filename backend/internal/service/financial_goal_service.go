package service

import (
	"context"
	"fmt"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/google/uuid"
)

type FinancialGoalService interface {
	CreateGoal(ctx context.Context, userID uuid.UUID, req model.CreateFinancialGoalRequest) (*model.FinancialGoal, error)
	ListGoals(ctx context.Context, userID uuid.UUID) ([]model.FinancialGoal, error)
	UpdateGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID, req model.UpdateFinancialGoalRequest) (*model.FinancialGoal, error)
	DeleteGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	LinkAsset(ctx context.Context, userID uuid.UUID, goalID uuid.UUID, assetID uuid.UUID) error
	UnlinkAsset(ctx context.Context, userID uuid.UUID, goalID uuid.UUID, assetID uuid.UUID) error
}

type financialGoalService struct {
	goalRepo  repository.FinancialGoalRepository
	assetRepo repository.AssetRepository
}

func NewFinancialGoalService(goalRepo repository.FinancialGoalRepository, assetRepo repository.AssetRepository) FinancialGoalService {
	return &financialGoalService{
		goalRepo:  goalRepo,
		assetRepo: assetRepo,
	}
}

func (s *financialGoalService) CreateGoal(ctx context.Context, userID uuid.UUID, req model.CreateFinancialGoalRequest) (*model.FinancialGoal, error) {
	goal := &model.FinancialGoal{
		UserID:       userID,
		Name:         req.Name,
		TargetAmount: req.TargetAmount,
		Priority:     req.Priority,
		TargetDate:   req.TargetDate,
		Status:       "active",
	}

	if err := s.goalRepo.Create(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}
	return goal, nil
}

func (s *financialGoalService) ListGoals(ctx context.Context, userID uuid.UUID) ([]model.FinancialGoal, error) {
	return s.goalRepo.ListByUser(ctx, userID)
}

func (s *financialGoalService) UpdateGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID, req model.UpdateFinancialGoalRequest) (*model.FinancialGoal, error) {
	goal, err := s.goalRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	if req.Name != "" {
		goal.Name = req.Name
	}
	if req.TargetAmount > 0 {
		goal.TargetAmount = req.TargetAmount
	}
	if req.Priority > 0 {
		goal.Priority = req.Priority
	}
	if req.TargetDate != nil {
		goal.TargetDate = req.TargetDate
	}
	if req.Status != "" {
		goal.Status = req.Status
	}

	if err := s.goalRepo.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return goal, nil
}

func (s *financialGoalService) DeleteGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.goalRepo.Delete(ctx, id, userID)
}

func (s *financialGoalService) LinkAsset(ctx context.Context, userID uuid.UUID, goalID uuid.UUID, assetID uuid.UUID) error {
	// Verify goal ownership
	if _, err := s.goalRepo.GetByID(ctx, goalID, userID); err != nil {
		return fmt.Errorf("goal not found or unauthorized: %w", err)
	}

	// Verify asset ownership
	asset, err := s.assetRepo.GetAssetByID(ctx, assetID, userID)
	if err != nil || asset == nil {
		return fmt.Errorf("asset not found or unauthorized: %w", err)
	}

	return s.goalRepo.LinkAsset(ctx, goalID, assetID)
}

func (s *financialGoalService) UnlinkAsset(ctx context.Context, userID uuid.UUID, goalID uuid.UUID, assetID uuid.UUID) error {
	// Verify goal ownership
	if _, err := s.goalRepo.GetByID(ctx, goalID, userID); err != nil {
		return fmt.Errorf("goal not found or unauthorized: %w", err)
	}

	return s.goalRepo.UnlinkAsset(ctx, goalID, assetID)
}
