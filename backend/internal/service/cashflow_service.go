package service

import (
	"context"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/google/uuid"
)

type CashflowService interface {
	CreateIncomeStream(ctx context.Context, userID uuid.UUID, req *model.CreateIncomeStreamRequest) (*model.IncomeStream, error)
	GetIncomeStreams(ctx context.Context, userID uuid.UUID) ([]*model.IncomeStream, error)
	UpdateIncomeStream(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateIncomeStreamRequest) (*model.IncomeStream, error)
	DeleteIncomeStream(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	CreateDependent(ctx context.Context, userID uuid.UUID, req *model.CreateDependentRequest) (*model.Dependent, error)
	GetDependents(ctx context.Context, userID uuid.UUID) ([]*model.Dependent, error)
	UpdateDependent(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateDependentRequest) (*model.Dependent, error)
	DeleteDependent(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type cashflowService struct {
	incomeRepo    *repository.IncomeRepository
	dependentRepo *repository.DependentRepository
}

func NewCashflowService(incomeRepo *repository.IncomeRepository, dependentRepo *repository.DependentRepository) CashflowService {
	return &cashflowService{
		incomeRepo:    incomeRepo,
		dependentRepo: dependentRepo,
	}
}

func (s *cashflowService) CreateIncomeStream(ctx context.Context, userID uuid.UUID, req *model.CreateIncomeStreamRequest) (*model.IncomeStream, error) {
	income := &model.IncomeStream{
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		IsPassive: req.IsPassive,
		Amount:    req.Amount,
		Frequency: req.Frequency,
		Notes:     req.Notes,
	}

	if err := s.incomeRepo.CreateIncomeStream(ctx, income); err != nil {
		return nil, err
	}
	return income, nil
}

func (s *cashflowService) GetIncomeStreams(ctx context.Context, userID uuid.UUID) ([]*model.IncomeStream, error) {
	return s.incomeRepo.GetIncomeStreamsByUserID(ctx, userID)
}

func (s *cashflowService) UpdateIncomeStream(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateIncomeStreamRequest) (*model.IncomeStream, error) {
	income := &model.IncomeStream{
		ID:        id,
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		IsPassive: req.IsPassive,
		Amount:    req.Amount,
		Frequency: req.Frequency,
		IsActive:  req.IsActive,
		Notes:     req.Notes,
	}

	if err := s.incomeRepo.UpdateIncomeStream(ctx, income); err != nil {
		return nil, err
	}
	return income, nil
}

func (s *cashflowService) DeleteIncomeStream(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.incomeRepo.DeleteIncomeStream(ctx, id, userID)
}

func (s *cashflowService) CreateDependent(ctx context.Context, userID uuid.UUID, req *model.CreateDependentRequest) (*model.Dependent, error) {
	dep := &model.Dependent{
		UserID:       userID,
		Name:         req.Name,
		Relationship: req.Relationship,
		DateOfBirth:  req.DateOfBirth,
		MonthlyCost:  req.MonthlyCost,
		Notes:        req.Notes,
	}

	if err := s.dependentRepo.CreateDependent(ctx, dep); err != nil {
		return nil, err
	}
	return dep, nil
}

func (s *cashflowService) GetDependents(ctx context.Context, userID uuid.UUID) ([]*model.Dependent, error) {
	return s.dependentRepo.GetDependentsByUserID(ctx, userID)
}

func (s *cashflowService) UpdateDependent(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateDependentRequest) (*model.Dependent, error) {
	dep := &model.Dependent{
		ID:           id,
		UserID:       userID,
		Name:         req.Name,
		Relationship: req.Relationship,
		DateOfBirth:  req.DateOfBirth,
		IsActive:     req.IsActive,
		MonthlyCost:  req.MonthlyCost,
		Notes:        req.Notes,
	}

	if err := s.dependentRepo.UpdateDependent(ctx, dep); err != nil {
		return nil, err
	}
	return dep, nil
}

func (s *cashflowService) DeleteDependent(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.dependentRepo.DeleteDependent(ctx, id, userID)
}
