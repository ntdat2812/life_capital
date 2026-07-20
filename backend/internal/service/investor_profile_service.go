package service

import (
	"context"
	"fmt"

	"time"

	"github.com/datnguyen/life_capital/backend/internal/ai"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/google/uuid"
)

type InvestorProfileService struct {
	repo        *repository.InvestorProfileRepository
	incomeRepo  *repository.IncomeRepository
	aiProviders []ai.AIProvider
}

func NewInvestorProfileService(repo *repository.InvestorProfileRepository, incomeRepo *repository.IncomeRepository, aiProviders ...ai.AIProvider) *InvestorProfileService {
	return &InvestorProfileService{
		repo:        repo,
		incomeRepo:  incomeRepo,
		aiProviders: aiProviders,
	}
}

func (s *InvestorProfileService) ProcessOnboarding(ctx context.Context, userID string, req *model.OnboardingRequest) (*model.InvestorProfile, error) {
	// 1. Format chat history into a single string
	chatText := ""
	for _, msg := range req.ChatHistory {
		chatText += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	// 2. Extract profile using AI with fallback mechanism
	extraction, err := ai.ExecuteWithFallback(s.aiProviders, func(p ai.AIProvider) (*ai.ExtractionResult, error) {
		return p.ExtractProfile(ctx, chatText)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to extract profile via AI: %w", err)
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// 3. Check for existing profile to determine version
	currentProfile, err := s.repo.GetActiveProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	newVersion := 1
	if currentProfile != nil {
		newVersion = currentProfile.Version + 1
		// Mark old profile as superseded
		if err := s.repo.SupersedePreviousProfiles(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to supersede old profiles: %w", err)
		}
	}

	// 4. Create new profile
	var dob *time.Time
	if extraction.DateOfBirth != nil {
		t, err := time.Parse("2006-01-02", *extraction.DateOfBirth)
		if err == nil {
			dob = &t
		}
	}

	profile := &model.InvestorProfile{
		UserID:                      uid,
		Version:                     newVersion,
		Status:                      "active",
		DateOfBirth:                 dob,
		MaritalStatus:               extraction.MaritalStatus,
		RiskTolerance:               extraction.RiskTolerance,
		RiskScore:                   extraction.RiskScore,
		EssentialMonthlyExpense:     extraction.EssentialMonthlyExpense,
		DiscretionaryMonthlyExpense: extraction.DiscretionaryMonthlyExpense,
		FITargetAmount:              extraction.FITargetAmount,
	}

	if err := s.repo.CreateProfile(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to save profile: %w", err)
	}

	// 5. If AI extracted an income, create an initial Income Stream
	if extraction.TotalMonthlyIncome > 0 {
		income := &model.IncomeStream{
			UserID:    uid,
			Name:      "Thu nhập chính (AI ước tính)",
			Type:      "salary",
			IsPassive: false,
			Amount:    extraction.TotalMonthlyIncome,
			Frequency: "monthly",
			IsActive:  true,
		}
		if err := s.incomeRepo.CreateIncomeStream(ctx, income); err != nil {
			// Just log the error, don't fail the onboarding
			fmt.Printf("Warning: failed to create initial income stream: %v\n", err)
		}
	}

	return profile, nil
}

func (s *InvestorProfileService) GetActiveProfile(ctx context.Context, userID string) (*model.InvestorProfile, error) {
	return s.repo.GetActiveProfileByUserID(ctx, userID)
}

func (s *InvestorProfileService) UpdateProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.InvestorProfile, error) {
	currentProfile, err := s.repo.GetActiveProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if currentProfile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	if req.DateOfBirth != nil {
		currentProfile.DateOfBirth = req.DateOfBirth
	}
	if req.MaritalStatus != "" {
		currentProfile.MaritalStatus = req.MaritalStatus
	}
	if req.RiskTolerance != "" {
		currentProfile.RiskTolerance = req.RiskTolerance
	}
	if req.RiskScore > 0 {
		currentProfile.RiskScore = req.RiskScore
	}
	if req.EssentialMonthlyExpense >= 0 {
		currentProfile.EssentialMonthlyExpense = req.EssentialMonthlyExpense
	}
	if req.DiscretionaryMonthlyExpense >= 0 {
		currentProfile.DiscretionaryMonthlyExpense = req.DiscretionaryMonthlyExpense
	}
	if req.FITargetAmount >= 0 {
		currentProfile.FITargetAmount = req.FITargetAmount
	}

	if err := s.repo.UpdateProfile(ctx, currentProfile); err != nil {
		return nil, err
	}

	return currentProfile, nil
}
