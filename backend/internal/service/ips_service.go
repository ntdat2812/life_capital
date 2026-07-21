package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/datnguyen/life_capital/backend/internal/ai"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/google/uuid"
)

type IPSService interface {
	GetLatestIPS(ctx context.Context, userID string) (*model.InvestmentPolicy, error)
	UpdateIPS(ctx context.Context, userID string, req *model.InvestmentPolicy) error
	GenerateIPS(ctx context.Context, userID string, preferredAssets []string) (*model.InvestmentPolicy, error)
}

type ipsService struct {
	ipsRepo      repository.InvestmentPolicyRepository
	profileRepo  *repository.InvestorProfileRepository
	assetRepo    repository.AssetRepository
	notifService NotificationService
	aiProviders  []ai.AIProvider
}

func NewIPSService(
	ipsRepo repository.InvestmentPolicyRepository,
	profileRepo *repository.InvestorProfileRepository,
	assetRepo repository.AssetRepository,
	notifService NotificationService,
	aiProviders []ai.AIProvider,
) IPSService {
	return &ipsService{
		ipsRepo:      ipsRepo,
		profileRepo:  profileRepo,
		assetRepo:    assetRepo,
		notifService: notifService,
		aiProviders:  aiProviders,
	}
}

func (s *ipsService) GetLatestIPS(ctx context.Context, userID string) (*model.InvestmentPolicy, error) {
	return s.ipsRepo.GetLatestByUserID(ctx, userID)
}

func (s *ipsService) UpdateIPS(ctx context.Context, userID string, req *model.InvestmentPolicy) error {
	policy, err := s.ipsRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if policy == nil {
		return fmt.Errorf("no IPS found to update")
	}

	policy.TargetAllocation = req.TargetAllocation
	// When user manually updates, it's no longer just an AI draft (if it was)
	policy.Status = "active"
	
	return s.ipsRepo.Update(ctx, policy)
}

func (s *ipsService) GenerateIPS(ctx context.Context, userID string, preferredAssets []string) (*model.InvestmentPolicy, error) {
	// 1. Fetch Profile
	profile, err := s.profileRepo.GetActiveProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active profile: %w", err)
	}
	if profile == nil {
		return nil, fmt.Errorf("no active profile found. please complete onboarding first")
	}

	// 2. Fetch Assets
	userUUID, err := uuid.Parse(userID)
	var assets []model.Asset = nil
	if err == nil {
		paginatedAssets, _ := s.assetRepo.GetAssetsByUserID(ctx, userUUID, "", "", 1000, 0)
		if paginatedAssets != nil {
			assets = paginatedAssets.Data
		}
	}

	var usedProvider string
	result, err := ai.ExecuteWithFallback(s.aiProviders, func(p ai.AIProvider) (*ai.IPSExtractionResult, error) {
		res, err := p.GenerateIPS(ctx, profile, assets, preferredAssets)
		if err == nil {
			usedProvider = p.Name()
		}
		return res, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI IPS: %w", err)
	}

	detailedStrategy := result.DetailedStrategy
	if usedProvider != "" {
		detailedStrategy += fmt.Sprintf("\n\n---\n*Phân tích và tạo tự động bởi AI Wealth Manager (Powered by **%s**)*", usedProvider)
	}

	targetAllocBytes, _ := json.Marshal(result.TargetAllocation)

	// 4. Update or Create IPS
	policy, err := s.ipsRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if policy == nil {
		// Create new
		policy = &model.InvestmentPolicy{
			UserID:           userID,
			Version:          1,
			ProfileVersion:   profile.Version,
			TargetAllocation: targetAllocBytes,
			DetailedStrategy: detailedStrategy,
			IsAIRecommended:  true,
			Status:           "active",
		}
		err = s.ipsRepo.Create(ctx, policy)
	} else {
		// Update existing
		policy.TargetAllocation = targetAllocBytes
		policy.DetailedStrategy = detailedStrategy
		policy.ProfileVersion = profile.Version
		policy.IsAIRecommended = true
		policy.Status = "active"
		err = s.ipsRepo.Update(ctx, policy)
	}

	if err != nil {
		return nil, err
	}

	// Send notification (Optional but good UX)
	actionLink := "/ips"
	_ = s.notifService.CreateNotification(
		context.Background(),
		userID,
		"IPS_RECOMMENDATION",
		"Chiến lược đầu tư đã được cập nhật",
		"AI đã phân tích và đưa ra bản tư vấn Chiến lược đầu tư mới nhất dành riêng cho bạn.",
		&actionLink,
	)

	return policy, nil
}
