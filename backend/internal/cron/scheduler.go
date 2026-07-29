package cron

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/datnguyen/life_capital/backend/internal/service"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron             *cron.Cron
	priceSyncService service.PriceSyncService
	userRepo         *repository.UserRepository
}

func NewScheduler(priceSyncService service.PriceSyncService, userRepo *repository.UserRepository) *Scheduler {
	c := cron.New()
	return &Scheduler{
		cron:             c,
		priceSyncService: priceSyncService,
		userRepo:         userRepo,
	}
}

func (s *Scheduler) Start() {
	schedule := os.Getenv("CRON_SCHEDULE_PRICE_SYNC")
	schedule = strings.Trim(schedule, `"`)
	if schedule == "" {
		schedule = "0 0 * * *" // default: everyday at midnight
	}

	_, err := s.cron.AddFunc(schedule, func() {
		log.Println("[CRON] Starting daily Price Sync for all users...")
		ctx := context.Background()
		
		// In a real system we would get all user IDs and sync them.
		// For our multi-user setup, we need a repository method to get all users.
		users, err := s.userRepo.GetAll(ctx)
		if err != nil {
			log.Printf("[CRON] Failed to get users for price sync: %v\n", err)
			return
		}

		for _, user := range users {
			report, err := s.priceSyncService.SyncByUser(ctx, user.ID)
			if err != nil {
				log.Printf("[CRON] Failed to sync prices for user %s: %v\n", user.ID, err)
			} else {
				log.Printf("[CRON] Price sync for user %s completed: %d updated, %d failed, %d skipped\n", user.ID, report.TotalUpdated, report.TotalFailed, report.TotalSkipped)
			}
		}
		
		log.Println("[CRON] Daily Price Sync completed.")
	})

	if err != nil {
		log.Fatalf("Failed to start Price Sync cron job: %v", err)
	}

	s.cron.Start()
	log.Printf("Cron scheduler started. Price Sync scheduled at: %s\n", schedule)
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
