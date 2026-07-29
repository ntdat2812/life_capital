package pricefeed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/datnguyen/life_capital/backend/internal/model"
)

type FmarketProvider struct{}

func NewFmarketProvider() *FmarketProvider {
	return &FmarketProvider{}
}

func (p *FmarketProvider) Name() string {
	return "Fmarket"
}

func (p *FmarketProvider) SupportedCategories() []model.AssetCategory {
	return []model.AssetCategory{model.AssetCategoryFund}
}

func (p *FmarketProvider) FetchPrices(ctx context.Context, tickers []string) (map[string]PriceResult, error) {
	results := make(map[string]PriceResult)
	if len(tickers) == 0 {
		return results, nil
	}

	url := "https://api.fmarket.vn/res/products/filter"
	payload := map[string]interface{}{
		"types":            []string{"NEW_FUND", "TRADING_FUND"},
		"isFmarketTrading": true,
		"searchField":      "",
		"sortOrder":        "DESC",
		"sortField":        "navToPrevious",
	}

	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Fmarket API returned status %d", resp.StatusCode)
	}

	var data struct {
		Data struct {
			Rows []struct {
				ShortName string  `json:"shortName"`
				Nav       float64 `json:"nav"`
			} `json:"rows"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Create a fast lookup map from the response
	navMap := make(map[string]float64)
	for _, row := range data.Data.Rows {
		navMap[row.ShortName] = row.Nav
	}

	// Match requested tickers
	for _, ticker := range tickers {
		if nav, exists := navMap[ticker]; exists {
			results[ticker] = PriceResult{
				Ticker:    ticker,
				Price:     nav,
				Source:    "Fmarket",
				FetchedAt: time.Now(),
			}
		}
	}

	return results, nil
}
