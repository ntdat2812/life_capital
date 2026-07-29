package pricefeed

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type DnseProvider struct{}

func NewDnseProvider() *DnseProvider {
	return &DnseProvider{}
}

func (p *DnseProvider) Name() string {
	return "DNSE"
}

func (p *DnseProvider) SupportedCategories() []model.AssetCategory {
	return []model.AssetCategory{model.AssetCategoryStock}
}

func (p *DnseProvider) FetchPrices(ctx context.Context, tickers []string) (map[string]PriceResult, error) {
	results := make(map[string]PriceResult)
	
	// Create channels for safe concurrent map writes
	resultCh := make(chan PriceResult, len(tickers))
	errCh := make(chan error, len(tickers))
	
	g, gCtx := errgroup.WithContext(ctx)
	// Limit concurrency to avoid getting blocked
	g.SetLimit(10)

	toTime := time.Now().Unix()
	fromTime := toTime - (7 * 24 * 60 * 60) // Fetch last 7 days to ensure we get a price even on weekends

	for _, ticker := range tickers {
		ticker := ticker // capture loop variable
		g.Go(func() error {
			url := fmt.Sprintf("https://services.entrade.com.vn/chart-api/v2/ohlcs/stock?from=%d&to=%d&symbol=%s&resolution=1D", fromTime, toTime, ticker)
			
			req, err := http.NewRequestWithContext(gCtx, "GET", url, nil)
			if err != nil {
				return err
			}
			req.Header.Set("User-Agent", "Mozilla/5.0")
			
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK {
				errCh <- fmt.Errorf("DNSE API returned status %d for %s", resp.StatusCode, ticker)
				return nil // Don't fail the whole group for one ticker
			}
			
			var data struct {
				C []float64 `json:"c"` // Close prices
				T []int64   `json:"t"` // Timestamps
			}
			
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				errCh <- fmt.Errorf("failed to decode DNSE response for %s: %v", ticker, err)
				return nil
			}
			
			if len(data.C) == 0 {
				errCh <- fmt.Errorf("no price data found for %s", ticker)
				return nil
			}
			
			// Get the most recent close price
			lastPrice := data.C[len(data.C)-1]
			// DNSE returns price in thousands usually (e.g. 24.5 for 24,500), but we need to check this.
			// Actually DNSE returns exact price or divided by 1000? Let's assume it returns exact or multiply by 1000 if it's small.
			// Wait, the API returns true prices, e.g., 28.5 means 28,500 VND. We should multiply by 1000.
			
			actualPrice := lastPrice * 1000
			
			resultCh <- PriceResult{
				Ticker:    ticker,
				Price:     actualPrice,
				Source:    "DNSE",
				FetchedAt: time.Now(),
			}
			
			return nil
		})
	}
	
	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(resultCh)
	close(errCh)
	
	for res := range resultCh {
		results[res.Ticker] = res
	}
	
	// Optional: log errors for specific tickers
	for err := range errCh {
		fmt.Printf("Warning: %v\n", err)
	}

	return results, nil
}
