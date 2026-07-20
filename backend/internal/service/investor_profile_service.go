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
	repo       *repository.InvestorProfileRepository
	aiProvider ai.AIProvider
}

func NewInvestorProfileService(repo *repository.InvestorProfileRepository, aiProvider ai.AIProvider) *InvestorProfileService {
	return &InvestorProfileService{
		repo:       repo,
		aiProvider: aiProvider,
	}
}

func (s *InvestorProfileService) ProcessOnboarding(ctx context.Context, userID string, req *model.OnboardingRequest) (*model.InvestorProfile, error) {
	// 1. Format chat history into a single string
	chatText := ""
	for _, msg := range req.ChatHistory {
		chatText += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	// 2. Extract profile using AI
	extraction, err := s.aiProvider.ExtractProfile(ctx, chatText)
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
		UserID:              uid,
		Version:             newVersion,
		Status:              "active",
		DateOfBirth:         dob,
		MaritalStatus:       extraction.MaritalStatus,
		RiskTolerance:       extraction.RiskTolerance,
		RiskScore:           extraction.RiskScore,
		TotalMonthlyIncome:  extraction.TotalMonthlyIncome,
		TotalMonthlyExpense: extraction.TotalMonthlyExpense,
		FITargetAmount:      extraction.FITargetAmount,
		LifeConstraints:     extraction.LifeConstraints,
	}

	if err := s.repo.CreateProfile(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to save profile: %w", err)
	}

	return profile, nil
}

func (s *InvestorProfileService) GetActiveProfile(ctx context.Context, userID string) (*model.InvestorProfile, error) {
	return s.repo.GetActiveProfileByUserID(ctx, userID)
}
