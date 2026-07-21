package service

import (
	"context"
	"fmt"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"

	"github.com/google/uuid"
)

type PortfolioService interface {
	// Portfolio (using Assets table)
	GetInvestableAssets(ctx context.Context, userID uuid.UUID) ([]model.Asset, error)

	// Thesis
	CreateThesis(ctx context.Context, thesis *model.InvestmentThesis) error
	GetThesesByUser(ctx context.Context, userID uuid.UUID) ([]*model.InvestmentThesis, error)
	GetThesisByTicker(ctx context.Context, userID uuid.UUID, ticker string) (*model.InvestmentThesis, error)
	UpdateThesis(ctx context.Context, thesis *model.InvestmentThesis) error
	DeleteThesis(ctx context.Context, id uuid.UUID) error

	// Watchlist
	CreateWatchlistItem(ctx context.Context, item *model.WatchlistItem) error
	GetWatchlistByUser(ctx context.Context, userID uuid.UUID) ([]*model.WatchlistItem, error)
	UpdateWatchlistItem(ctx context.Context, item *model.WatchlistItem) error
	DeleteWatchlistItem(ctx context.Context, id uuid.UUID) error
}

type portfolioService struct {
	assetRepo     repository.AssetRepository
	thesisRepo    repository.ThesisRepository
	watchlistRepo repository.WatchlistRepository
	txManager     *repository.TxManager
}

func NewPortfolioService(
	assetRepo repository.AssetRepository,
	thesisRepo repository.ThesisRepository,
	watchlistRepo repository.WatchlistRepository,
	txManager *repository.TxManager,
) PortfolioService {
	return &portfolioService{
		assetRepo:     assetRepo,
		thesisRepo:    thesisRepo,
		watchlistRepo: watchlistRepo,
		txManager:     txManager,
	}
}

func (s *portfolioService) GetInvestableAssets(ctx context.Context, userID uuid.UUID) ([]model.Asset, error) {
	assets, err := s.assetRepo.GetAssetsByUserID(ctx, userID, "", "", 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}

	var portfolio []model.Asset
	for _, asset := range assets.Data {
		// Filter only investable categories
		if asset.Category == "stock" || asset.Category == "crypto" || asset.Category == "real_estate" || asset.Category == "gold" || asset.Category == "fund" {
			portfolio = append(portfolio, asset)
		}
	}
	return portfolio, nil
}

func (s *portfolioService) CreateThesis(ctx context.Context, thesis *model.InvestmentThesis) error {
	return s.thesisRepo.Create(ctx, thesis)
}

func (s *portfolioService) GetThesesByUser(ctx context.Context, userID uuid.UUID) ([]*model.InvestmentThesis, error) {
	return s.thesisRepo.GetByUserID(ctx, userID)
}

func (s *portfolioService) GetThesisByTicker(ctx context.Context, userID uuid.UUID, ticker string) (*model.InvestmentThesis, error) {
	return s.thesisRepo.GetByTicker(ctx, userID, ticker)
}

func (s *portfolioService) UpdateThesis(ctx context.Context, thesis *model.InvestmentThesis) error {
	return s.thesisRepo.Update(ctx, thesis)
}

func (s *portfolioService) DeleteThesis(ctx context.Context, id uuid.UUID) error {
	return s.thesisRepo.Delete(ctx, id)
}

func (s *portfolioService) CreateWatchlistItem(ctx context.Context, item *model.WatchlistItem) error {
	return s.watchlistRepo.Create(ctx, item)
}

func (s *portfolioService) GetWatchlistByUser(ctx context.Context, userID uuid.UUID) ([]*model.WatchlistItem, error) {
	return s.watchlistRepo.GetByUserID(ctx, userID)
}

func (s *portfolioService) UpdateWatchlistItem(ctx context.Context, item *model.WatchlistItem) error {
	return s.watchlistRepo.Update(ctx, item)
}

func (s *portfolioService) DeleteWatchlistItem(ctx context.Context, id uuid.UUID) error {
	return s.watchlistRepo.Delete(ctx, id)
}
