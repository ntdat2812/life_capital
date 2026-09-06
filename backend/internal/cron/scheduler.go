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
		
		report, err := s.priceSyncService.SyncAll(ctx)
		if err != nil {
			log.Printf("[CRON] Failed to sync prices: %v\n", err)
			return
		}

		log.Printf("[CRON] Daily Price Sync completed: %d updated, %d failed, %d skipped\n", report.TotalUpdated, report.TotalFailed, report.TotalSkipped)
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
