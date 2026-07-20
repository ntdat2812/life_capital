package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
)

type WealthService interface {
	CreateAsset(ctx context.Context, userID uuid.UUID, req *model.CreateAssetRequest) (*model.Asset, error)
	GetAssets(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedAssets, error)
	UpdateAsset(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateAssetRequest) (*model.Asset, error)
	DeleteAsset(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	CreateLiability(ctx context.Context, userID uuid.UUID, req *model.CreateLiabilityRequest) (*model.Liability, error)
	GetLiabilities(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedLiabilities, error)
	UpdateLiability(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateLiabilityRequest) (*model.Liability, error)
	DeleteLiability(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	GetNetWorthSummary(ctx context.Context, userID uuid.UUID) (*model.NetWorthSummary, error)
}

type wealthService struct {
	assetRepo     repository.AssetRepository
	liabilityRepo repository.LiabilityRepository
	userRepo      *repository.UserRepository
}

func NewWealthService(assetRepo repository.AssetRepository, liabilityRepo repository.LiabilityRepository, userRepo *repository.UserRepository) WealthService {
	return &wealthService{
		assetRepo:     assetRepo,
		liabilityRepo: liabilityRepo,
		userRepo:      userRepo,
	}
}

func (s *wealthService) CreateAsset(ctx context.Context, userID uuid.UUID, req *model.CreateAssetRequest) (*model.Asset, error) {
	asset := &model.Asset{
		UserID:       userID,
		Category:     req.Category,
		Name:         req.Name,
		Ticker:       req.Ticker,
		Quantity:     req.Quantity,
		AvgPrice:     req.AvgPrice,
		CurrentPrice: req.CurrentPrice,
		CurrentValue: req.CurrentValue,
		CostBasis:    req.CostBasis,
		Notes:        req.Notes,
	}

	if err := s.assetRepo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *wealthService) GetAssets(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedAssets, error) {
	return s.assetRepo.GetAssetsByUserID(ctx, userID, category, sort, limit, offset)
}

func (s *wealthService) UpdateAsset(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateAssetRequest) (*model.Asset, error) {
	asset := &model.Asset{
		ID:           id,
		UserID:       userID,
		Category:     req.Category,
		Name:         req.Name,
		Ticker:       req.Ticker,
		Quantity:     req.Quantity,
		AvgPrice:     req.AvgPrice,
		CurrentPrice: req.CurrentPrice,
		CurrentValue: req.CurrentValue,
		CostBasis:    req.CostBasis,
		Notes:        req.Notes,
	}

	if err := s.assetRepo.UpdateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *wealthService) DeleteAsset(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.assetRepo.DeleteAsset(ctx, id, userID)
}

func (s *wealthService) CreateLiability(ctx context.Context, userID uuid.UUID, req *model.CreateLiabilityRequest) (*model.Liability, error) {
	liability := &model.Liability{
		UserID:           userID,
		Category:         req.Category,
		Name:             req.Name,
		RemainingBalance: req.RemainingBalance,
		InterestRate:     req.InterestRate,
		MonthlyPayment:   req.MonthlyPayment,
		Lender:           req.Lender,
		Notes:            req.Notes,
	}

	if err := s.liabilityRepo.CreateLiability(ctx, liability); err != nil {
		return nil, err
	}

	return liability, nil
}

func (s *wealthService) GetLiabilities(ctx context.Context, userID uuid.UUID, category string, sort string, limit int, offset int) (*model.PaginatedLiabilities, error) {
	return s.liabilityRepo.GetLiabilitiesByUserID(ctx, userID, category, sort, limit, offset)
}

func (s *wealthService) UpdateLiability(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.UpdateLiabilityRequest) (*model.Liability, error) {
	liability := &model.Liability{
		ID:               id,
		UserID:           userID,
		Category:         req.Category,
		Name:             req.Name,
		RemainingBalance: req.RemainingBalance,
		InterestRate:     req.InterestRate,
		MonthlyPayment:   req.MonthlyPayment,
		Lender:           req.Lender,
		Notes:            req.Notes,
	}

	if err := s.liabilityRepo.UpdateLiability(ctx, liability); err != nil {
		return nil, err
	}

	return liability, nil
}

func (s *wealthService) DeleteLiability(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.liabilityRepo.DeleteLiability(ctx, id, userID)
}

func (s *wealthService) GetNetWorthSummary(ctx context.Context, userID uuid.UUID) (*model.NetWorthSummary, error) {
	// For net worth summary, we need ALL assets without limit (limit=0, offset=0)
	assetsPage, err := s.assetRepo.GetAssetsByUserID(ctx, userID, "", "", 0, 0)
	if err != nil {
		return nil, err
	}
	assets := assetsPage.Data

	liabilitiesPage, err := s.liabilityRepo.GetLiabilitiesByUserID(ctx, userID, "", "", 0, 0)
	if err != nil {
		return nil, err
	}
	liabilities := liabilitiesPage.Data

	user, err := s.userRepo.FindByID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	var totalAssets float64
	for _, a := range assets {
		totalAssets += a.CurrentValue
	}

	var totalLiabilities float64
	for _, l := range liabilities {
		totalLiabilities += l.RemainingBalance
	}

	baseCurrency := "VND"
	if user != nil && user.BaseCurrency != "" {
		baseCurrency = user.BaseCurrency
	}

	return &model.NetWorthSummary{
		TotalAssets:      totalAssets,
		TotalLiabilities: totalLiabilities,
		NetWorth:         totalAssets - totalLiabilities,
		BaseCurrency:     baseCurrency,
	}, nil
}
