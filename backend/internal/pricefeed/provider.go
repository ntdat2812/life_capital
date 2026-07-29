package pricefeed

import (
	"context"
	"time"

	"github.com/datnguyen/life_capital/backend/internal/model"
)

type PriceResult struct {
	Ticker    string
	Price     float64
	Source    string
	FetchedAt time.Time
}

type PriceProvider interface {
	Name() string
	SupportedCategories() []model.AssetCategory
	FetchPrices(ctx context.Context, tickers []string) (map[string]PriceResult, error)
}

type ProviderRegistry struct {
	providers map[model.AssetCategory]PriceProvider
}

func NewRegistry(providers ...PriceProvider) *ProviderRegistry {
	reg := &ProviderRegistry{
		providers: make(map[model.AssetCategory]PriceProvider),
	}

	for _, p := range providers {
		for _, cat := range p.SupportedCategories() {
			reg.providers[cat] = p
		}
	}
	return reg
}

func (r *ProviderRegistry) GetProvider(cat model.AssetCategory) (PriceProvider, bool) {
	p, exists := r.providers[cat]
	return p, exists
}
