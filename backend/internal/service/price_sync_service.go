package service

import (
	"context"
	"fmt"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/pricefeed"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type SyncReport struct {
	TotalUpdated int                      `json:"total_updated"`
	TotalFailed  int                      `json:"total_failed"`
	TotalSkipped int                      `json:"total_skipped"`
	Details      []map[string]interface{} `json:"details"`
}

type PriceSyncService interface {
	SyncByUser(ctx context.Context, userID uuid.UUID) (*SyncReport, error)
}

type priceSyncService struct {
	assetRepo repository.AssetRepository
	registry  *pricefeed.ProviderRegistry
}

func NewPriceSyncService(assetRepo repository.AssetRepository, registry *pricefeed.ProviderRegistry) PriceSyncService {
	return &priceSyncService{
		assetRepo: assetRepo,
		registry:  registry,
	}
}

func (s *priceSyncService) SyncByUser(ctx context.Context, userID uuid.UUID) (*SyncReport, error) {
	// 1. Get all syncable assets (have ticker and quantity)
	assets, err := s.assetRepo.GetSyncableAssets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get syncable assets: %w", err)
	}

	if len(assets) == 0 {
		return &SyncReport{}, nil
	}

	// 2. Group assets by category and collect tickers
	categoryTickers := make(map[model.AssetCategory][]string)
	assetMap := make(map[string][]model.Asset) // Ticker -> []Asset (in case of multiple assets with same ticker)

	for _, a := range assets {
		if a.Ticker == nil {
			continue
		}
		ticker := *a.Ticker
		categoryTickers[a.Category] = append(categoryTickers[a.Category], ticker)
		assetMap[ticker] = append(assetMap[ticker], a)
	}

	// 3. Fetch prices concurrently for each category
	priceResults := make(map[string]pricefeed.PriceResult)
	var report SyncReport

	g, gCtx := errgroup.WithContext(ctx)
	resultsCh := make(chan map[string]pricefeed.PriceResult, len(categoryTickers))

	for cat, tickers := range categoryTickers {
		cat := cat
		tickers := uniqueStrings(tickers)
		
		provider, exists := s.registry.GetProvider(cat)
		if !exists {
			// No provider for this category, skip
			report.TotalSkipped += len(tickers)
			continue
		}

		g.Go(func() error {
			res, err := provider.FetchPrices(gCtx, tickers)
			if err != nil {
				// We don't fail the whole sync if one provider fails, just log and return nil
				fmt.Printf("Provider %s failed: %v\n", provider.Name(), err)
				return nil 
			}
			resultsCh <- res
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(resultsCh)

	// Merge results
	for resMap := range resultsCh {
		for ticker, res := range resMap {
			priceResults[ticker] = res
		}
	}

	// 4. Prepare updates
	var updates []model.PriceUpdate
	for ticker, priceResult := range priceResults {
		if priceResult.Price <= 0 {
			continue
		}

		matchedAssets := assetMap[ticker]
		for _, a := range matchedAssets {
			qty := *a.Quantity
			currentValue := qty * priceResult.Price

			updates = append(updates, model.PriceUpdate{
				ID:           a.ID,
				CurrentPrice: priceResult.Price,
				CurrentValue: currentValue,
			})

			report.TotalUpdated++
			report.Details = append(report.Details, map[string]interface{}{
				"id":     a.ID,
				"name":   a.Name,
				"ticker": ticker,
				"price":  priceResult.Price,
				"source": priceResult.Source,
				"status": "success",
			})
		}
	}

	// Count failed
	report.TotalFailed = len(assets) - report.TotalUpdated - report.TotalSkipped

	// 5. Batch update DB
	if len(updates) > 0 {
		if err := s.assetRepo.BatchUpdatePrices(ctx, updates); err != nil {
			return nil, fmt.Errorf("failed to batch update prices: %w", err)
		}
	}

	return &report, nil
}

func uniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
