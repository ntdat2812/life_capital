package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/datnguyen/life_capital/backend/internal/ai"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TimelineService struct {
	aiProviders  []ai.AIProvider
	profileRepo  *repository.InvestorProfileRepository
	incomeRepo   *repository.IncomeRepository
	depRepo      *repository.DependentRepository
	eventRepo    repository.LifeEventRepository
	txManager    *repository.TxManager
}

func NewTimelineService(
	aiProviders []ai.AIProvider,
	profileRepo *repository.InvestorProfileRepository,
	incomeRepo *repository.IncomeRepository,
	depRepo *repository.DependentRepository,
	eventRepo repository.LifeEventRepository,
	txManager *repository.TxManager,
) *TimelineService {
	return &TimelineService{
		aiProviders:  aiProviders,
		profileRepo:  profileRepo,
		incomeRepo:   incomeRepo,
		depRepo:      depRepo,
		eventRepo:    eventRepo,
		txManager:    txManager,
	}
}

func (s *TimelineService) AnalyzeEvent(ctx context.Context, userID, eventText string) (*ai.LifeEventAnalysisResult, error) {
	var (
		profile *model.InvestorProfile
		incomes []*model.IncomeStream
		deps    []*model.Dependent
		errProf error
		errInc  error
		errDep  error
		wg      sync.WaitGroup
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		profile, errProf = s.profileRepo.GetActiveProfileByUserID(ctx, userID)
	}()

	go func() {
		defer wg.Done()
		uid, err := uuid.Parse(userID)
		if err == nil {
			incomes, errInc = s.incomeRepo.GetIncomeStreamsByUserID(ctx, uid)
		} else {
			errInc = err
		}
	}()

	go func() {
		defer wg.Done()
		uid, err := uuid.Parse(userID)
		if err == nil {
			deps, errDep = s.depRepo.GetDependentsByUserID(ctx, uid)
		} else {
			errDep = err
		}
	}()

	wg.Wait()

	if errProf != nil {
		return nil, fmt.Errorf("failed to load profile: %w", errProf)
	}
	if errInc != nil {
		return nil, fmt.Errorf("failed to load incomes: %w", errInc)
	}
	if errDep != nil {
		return nil, fmt.Errorf("failed to load dependents: %w", errDep)
	}

	// Format prompt context
	var sb strings.Builder

	if profile != nil {
		sb.WriteString(fmt.Sprintf("Hồ sơ rủi ro hiện tại (Risk Score): %d (Khẩu vị: %s)\n", profile.RiskScore, profile.RiskTolerance))
		sb.WriteString(fmt.Sprintf("Tình trạng hôn nhân: %s\n", profile.MaritalStatus))
		// We need total income. Let's compute it.
		var totalIncome float64
		for _, inc := range incomes {
			if inc.IsActive {
				totalIncome += inc.Amount
			}
		}
		sb.WriteString(fmt.Sprintf("Thu nhập hàng tháng: %.0f\n", totalIncome))
		sb.WriteString(fmt.Sprintf("Chi phí thiết yếu hiện tại: %.0f\n", profile.EssentialMonthlyExpense))
		sb.WriteString(fmt.Sprintf("Chi phí hưởng thụ hiện tại: %.0f\n", profile.DiscretionaryMonthlyExpense))
	} else {
		sb.WriteString("Hồ sơ rủi ro hiện tại: Chưa có\n")
	}

	sb.WriteString("Các nguồn thu nhập hiện tại:\n")
	if len(incomes) == 0 {
		sb.WriteString("- Không có\n")
	} else {
		for _, inc := range incomes {
			if inc.IsActive {
				sb.WriteString(fmt.Sprintf("- %s (%s): %.0f\n", inc.Name, inc.Type, inc.Amount))
			}
		}
	}

	sb.WriteString("Người phụ thuộc hiện tại:\n")
	if len(deps) == 0 {
		sb.WriteString("- Không có\n")
	} else {
		for _, d := range deps {
			if d.IsActive {
				sb.WriteString(fmt.Sprintf("- %s (%s): %.0f\n", d.Name, d.Relationship, d.MonthlyCost))
			}
		}
	}

	sb.WriteString("\n--- SỰ KIỆN MỚI ---\n")
	sb.WriteString(fmt.Sprintf("\"%s\"\n", eventText))

	// Call AI Provider with fallback logic using ExecuteWithFallback wrapper?
	// Oh wait, AIProvider has AnalyzeLifeEvent. But we should wrap it with ExecuteWithFallback.
	// `ai.ExecuteWithFallback` calls a function for each provider. Let's see ai_provider.go.
	return ai.ExecuteWithFallback(s.aiProviders, func(p ai.AIProvider) (*ai.LifeEventAnalysisResult, error) {
		return p.AnalyzeLifeEvent(ctx, sb.String())
	})
}

func (s *TimelineService) ConfirmEvent(ctx context.Context, userID string, req *model.ConfirmEventRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return s.txManager.ExecTx(ctx, func(tx pgx.Tx) error {
		txProfileRepo := repository.NewInvestorProfileRepository(tx)
		txEventRepo := repository.NewLifeEventRepository(tx)
		txIncomeRepo := repository.NewIncomeRepository(tx)
		txDepRepo := repository.NewDependentRepository(tx)

		// 1. Get current profile
		currentProfile, err := txProfileRepo.GetActiveProfileByUserID(ctx, userID)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}

		newVersion := 1
		
		now := time.Now()
		newProfile := &model.InvestorProfile{
			UserID:                      uid,
			Status:                      "active",
			DateOfBirth:                 &now, // fallback
			MaritalStatus:               "Độc thân",
			RiskTolerance:               "Trung bình",
			RiskScore:                   50,
			EssentialMonthlyExpense:     0,
			DiscretionaryMonthlyExpense: 0,
			FITargetAmount:              0,
		}

		if currentProfile != nil {
			newVersion = currentProfile.Version + 1
			
			// Copy values
			newProfile.DateOfBirth = currentProfile.DateOfBirth
			newProfile.MaritalStatus = currentProfile.MaritalStatus
			newProfile.RiskTolerance = currentProfile.RiskTolerance
			newProfile.RiskScore = currentProfile.RiskScore
			newProfile.EssentialMonthlyExpense = currentProfile.EssentialMonthlyExpense
			newProfile.DiscretionaryMonthlyExpense = currentProfile.DiscretionaryMonthlyExpense
			newProfile.FITargetAmount = currentProfile.FITargetAmount

			// Supersede old profile
			if err := txProfileRepo.SupersedePreviousProfiles(ctx, userID); err != nil {
				return err
			}
		}

		newProfile.Version = newVersion

		// Apply impacts
		if req.NewMaritalStatus != nil {
			newProfile.MaritalStatus = *req.NewMaritalStatus
		}
		if req.NewRiskScore != nil {
			newProfile.RiskScore = *req.NewRiskScore
		}
		if req.NewRiskTolerance != nil {
			newProfile.RiskTolerance = *req.NewRiskTolerance
		}
		newProfile.EssentialMonthlyExpense += req.ExpenseImpact

		// 2. Create new profile
		if err := txProfileRepo.CreateProfile(ctx, newProfile); err != nil {
			return err
		}

		// 3. Create life event
		event := &model.LifeEvent{
			UserID:                  uid,
			EventDate:               time.Now(),
			Category:                req.Category,
			Title:                   req.Title,
			IncomeImpact:            req.IncomeImpact,
			ExpenseImpact:           req.ExpenseImpact,
			AIImpactAnalysis:        &req.AIImpactAnalysis,
			TriggeredProfileVersion: &newVersion,
			RequiresIPSUpdate:       false,
		}

		if err := txEventRepo.Create(ctx, event); err != nil {
			return err
		}

		// 4. Add new income streams
		for _, inc := range req.IncomeStreamsToAdd {
			err := txIncomeRepo.CreateIncomeStream(ctx, &model.IncomeStream{
				UserID:    uid,
				Name:      inc.Name,
				Type:      inc.Type,
				Amount:    inc.Amount,
				IsPassive: inc.IsPassive,
				IsActive:  true,
			})
			if err != nil {
				return err
			}
		}

		// 5. Add new dependents
		for _, dep := range req.DependentsToAdd {
			err := txDepRepo.CreateDependent(ctx, &model.Dependent{
				UserID:       uid,
				Name:         dep.Name,
				Relationship: dep.Relationship,
				MonthlyCost:  dep.MonthlyCost,
				IsActive:     true,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *TimelineService) GetTimeline(ctx context.Context, userID string) ([]model.LifeEvent, error) {
	return s.eventRepo.GetByUserID(ctx, userID)
}
