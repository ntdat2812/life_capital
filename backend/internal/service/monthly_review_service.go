package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"github.com/google/uuid"
	"github.com/datnguyen/life_capital/backend/internal/ai"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
)

type MonthlyReviewService interface {
	GenerateReview(ctx context.Context, userID uuid.UUID, newInvestmentAmount float64) (*model.MonthlyReview, error)
	GetReviewHistory(ctx context.Context, userID uuid.UUID) ([]model.MonthlyReview, error)
	GetReviewByMonth(ctx context.Context, userID uuid.UUID, month string) (*model.MonthlyReview, error)
	SaveReview(ctx context.Context, review *model.MonthlyReview) error
}

type monthlyReviewService struct {
	monthlyReviewRepo repository.MonthlyReviewRepository
	portfolioSvc      PortfolioService
	liabilityRepo     repository.LiabilityRepository
	ipsRepo           repository.InvestmentPolicyRepository
	assetRepo         repository.AssetRepository
	dependentRepo     *repository.DependentRepository
	thesisRepo        repository.ThesisRepository
	aiProviders       []ai.AIProvider
}

func NewMonthlyReviewService(
	monthlyReviewRepo repository.MonthlyReviewRepository,
	portfolioSvc PortfolioService,
	liabilityRepo repository.LiabilityRepository,
	ipsRepo repository.InvestmentPolicyRepository,
	assetRepo repository.AssetRepository,
	dependentRepo *repository.DependentRepository,
	thesisRepo repository.ThesisRepository,
	aiProviders []ai.AIProvider,
) MonthlyReviewService {
	return &monthlyReviewService{
		monthlyReviewRepo: monthlyReviewRepo,
		portfolioSvc:      portfolioSvc,
		liabilityRepo:     liabilityRepo,
		ipsRepo:           ipsRepo,
		assetRepo:         assetRepo,
		dependentRepo:     dependentRepo,
		thesisRepo:        thesisRepo,
		aiProviders:       aiProviders,
	}
}

func (s *monthlyReviewService) GetReviewHistory(ctx context.Context, userID uuid.UUID) ([]model.MonthlyReview, error) {
	return s.monthlyReviewRepo.ListByUser(ctx, userID)
}

func (s *monthlyReviewService) GetReviewByMonth(ctx context.Context, userID uuid.UUID, month string) (*model.MonthlyReview, error) {
	return s.monthlyReviewRepo.GetByMonth(ctx, userID, month)
}

func (s *monthlyReviewService) GenerateReview(ctx context.Context, userID uuid.UUID, newInvestmentAmount float64) (*model.MonthlyReview, error) {
	var ips *model.InvestmentPolicy
	var assetsResp *model.PaginatedAssets
	var watchlist []*model.WatchlistItem
	var liabilitiesResp *model.PaginatedLiabilities
	var dependents []*model.Dependent
	var theses []*model.InvestmentThesis

	g, gCtx := errgroup.WithContext(ctx)

	// 1. Get IPS
	g.Go(func() error {
		var err error
		ips, err = s.ipsRepo.GetLatestByUserID(gCtx, userID.String())
		if err != nil {
			return fmt.Errorf("failed to get active IPS: %v", err)
		}
		return nil
	})

	// 2. Get Portfolio (Assets) - ALL assets
	g.Go(func() error {
		var err error
		assetsResp, err = s.assetRepo.GetAssetsByUserID(gCtx, userID, "", "", 100, 0)
		if err != nil {
			return fmt.Errorf("failed to get assets: %v", err)
		}
		return nil
	})

	// 3. Get Watchlist
	g.Go(func() error {
		var err error
		watchlist, err = s.portfolioSvc.GetWatchlistByUser(gCtx, userID)
		if err != nil {
			return fmt.Errorf("failed to get watchlist: %v", err)
		}
		return nil
	})

	// 4. Get Liabilities
	g.Go(func() error {
		var err error
		liabilitiesResp, err = s.liabilityRepo.GetLiabilitiesByUserID(gCtx, userID, "", "", 100, 0)
		if err != nil {
			return fmt.Errorf("failed to get liabilities: %v", err)
		}
		return nil
	})

	// 5. Get Dependents
	g.Go(func() error {
		var err error
		dependents, err = s.dependentRepo.GetDependentsByUserID(gCtx, userID)
		if err != nil {
			return fmt.Errorf("failed to get dependents: %v", err)
		}
		return nil
	})

	// 6. Get Investment Theses
	g.Go(func() error {
		var err error
		theses, err = s.thesisRepo.GetByUserID(gCtx, userID)
		if err != nil {
			return fmt.Errorf("failed to get theses: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	assets := assetsResp.Data

	// Calculate current net worth inside portfolio
	var portfolioValue float64
	for _, a := range assets {
		portfolioValue += a.CurrentValue
	}

	var totalDebt float64
	for _, l := range liabilitiesResp.Data {
		totalDebt += l.RemainingBalance
	}

	netWorth := portfolioValue - totalDebt // Simplified net worth for review context

	// Prepare JSON strings for AI Context
	ipsJSON, _ := json.MarshalIndent(ips.TargetAllocation, "", "  ")

	// Simplify assets for prompt to save tokens
	type simplifiedAsset struct {
		Category    string  `json:"category"`
		Name        string  `json:"name"`
		Ticker      string  `json:"ticker"`
		CostBasis   float64 `json:"cost_basis"`
		Value       float64 `json:"value"`
		Profit      float64 `json:"profit"`
		ROI         float64 `json:"roi_percentage"`
	}
	var simpAssets []simplifiedAsset
	for _, a := range assets {
		ticker := ""
		if a.Ticker != nil {
			ticker = *a.Ticker
		}
		costBasis := 0.0
		if a.CostBasis != nil {
			costBasis = *a.CostBasis
		}
		profit := a.CurrentValue - costBasis
		roi := 0.0
		if costBasis > 0 {
			roi = (profit / costBasis) * 100
		}
		
		simpAssets = append(simpAssets, simplifiedAsset{
			Category:  string(a.Category),
			Name:      a.Name,
			Ticker:    ticker,
			CostBasis: costBasis,
			Value:     a.CurrentValue,
			Profit:    profit,
			ROI:       roi,
		})
	}
	assetsJSON, _ := json.MarshalIndent(simpAssets, "", "  ")
	dependentsJSON, _ := json.MarshalIndent(dependents, "", "  ")

	type simplifiedWatchlist struct {
		Ticker      string  `json:"ticker"`
		TargetPrice float64 `json:"target_price"`
		Priority    int     `json:"priority"`
	}
	var simpWatchlist []simplifiedWatchlist
	for _, w := range watchlist {
		simpWatchlist = append(simpWatchlist, simplifiedWatchlist{
			Ticker:      w.Ticker,
			TargetPrice: w.TargetPrice,
			Priority:    w.Priority,
		})
	}
	watchlistJSON, _ := json.MarshalIndent(simpWatchlist, "", "  ")

	type simplifiedLiability struct {
		Category     string  `json:"category"`
		Name         string  `json:"name"`
		Balance      float64 `json:"balance"`
		InterestRate float64 `json:"interest_rate"`
	}
	var simpLiabs []simplifiedLiability
	for _, l := range liabilitiesResp.Data {
		interestRate := 0.0
		if l.InterestRate != nil {
			interestRate = *l.InterestRate
		}
		simpLiabs = append(simpLiabs, simplifiedLiability{
			Category:     string(l.Category),
			Name:         l.Name,
			Balance:      l.RemainingBalance,
			InterestRate: interestRate,
		})
	}
	liabilitiesJSON, _ := json.MarshalIndent(simpLiabs, "", "  ")

	type simplifiedThesis struct {
		Ticker        string          `json:"ticker"`
		CompanyName   string          `json:"company_name"`
		ThesisSummary string          `json:"thesis_summary"`
		Catalysts     json.RawMessage `json:"catalysts"`
		Risks         json.RawMessage `json:"risks"`
		Notes         string          `json:"notes"`
	}
	var simpTheses []simplifiedThesis
	for _, t := range theses {
		simpTheses = append(simpTheses, simplifiedThesis{
			Ticker:        t.Ticker,
			CompanyName:   t.CompanyName,
			ThesisSummary: t.ThesisSummary,
			Catalysts:     t.Catalysts,
			Risks:         t.Risks,
			Notes:         t.Notes,
		})
	}
	thesesJSON, _ := json.MarshalIndent(simpTheses, "", "  ")

	// Build template replacements for the prompt
	reviewReplacements := map[string]string{
		"{{.IPS_JSON}}":          string(ipsJSON),
		"{{.HOLDINGS_JSON}}":     string(assetsJSON),
		"{{.WATCHLIST_JSON}}":    string(watchlistJSON),
		"{{.LIABILITIES_JSON}}":  string(liabilitiesJSON),
		"{{.DEPENDENTS_JSON}}":   string(dependentsJSON),
		"{{.THESES_JSON}}":       string(thesesJSON),
		"{{.NEW_CAPITAL}}":       fmt.Sprintf("%.0f", newInvestmentAmount),
	}

	// Call AI Provider
	var usedAIProvider string
	aiResp, err := ai.ExecuteWithFallback(s.aiProviders, func(p ai.AIProvider) (*model.MonthlyReviewRecommendationResponse, error) {
		usedAIProvider = p.Name() + " (" + p.Model() + ")"
		return p.GenerateMonthlyReview(ctx, reviewReplacements)
	})
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %v", err)
	}

	if aiResp != nil {
		aiResp.AIOverallSummary = fmt.Sprintf("%s\n\n*(Đánh giá bởi AI: %s)*", aiResp.AIOverallSummary, usedAIProvider)
	}

	now := time.Now()
	monthStr := fmt.Sprintf("%04d-%02d-01", now.Year(), now.Month()) // Store as first day of the month
	reviewMonthTime, _ := time.Parse("2006-01-02", monthStr)

	// Return a PREVIEW object, do not save to DB yet.
	review := &model.MonthlyReview{
		UserID:              userID,
		ReviewMonth:         reviewMonthTime,
		Status:              "draft", // Indicates it's a preview
		NewInvestmentAmount: newInvestmentAmount,
		PortfolioSnapshot:   simpAssets, // Saving simplified snapshot
		NetWorthAtReview:    netWorth,
		AIRecommendations:   aiResp.AIRecommendations,
		AIOverallSummary:    aiResp.AIOverallSummary,
	}

	return review, nil
}

func (s *monthlyReviewService) SaveReview(ctx context.Context, review *model.MonthlyReview) error {
	// Set status to completed since user confirmed it
	review.Status = "completed"

	// Check if already exists for this month
	monthStr := fmt.Sprintf("%04d-%02d-01", review.ReviewMonth.Year(), review.ReviewMonth.Month())
	existing, err := s.monthlyReviewRepo.GetByMonth(ctx, review.UserID, monthStr)

	if err == nil && existing != nil {
		// Update existing
		review.ID = existing.ID
		return s.monthlyReviewRepo.Update(ctx, review)
	}

	// Create new
	return s.monthlyReviewRepo.Create(ctx, review)
}
