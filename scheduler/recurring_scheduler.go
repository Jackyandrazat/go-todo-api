package scheduler

import (
	"log"
	"time"

	"go-todo-api/config"
	"go-todo-api/model"
	"go-todo-api/repository"

	"github.com/robfig/cron/v3"
)

type RecurringScheduler struct {
	repo *repository.RecurringTransactionRepository
}

func NewRecurringScheduler() *RecurringScheduler {
	return &RecurringScheduler{
		repo: repository.NewRecurringTransactionRepository(),
	}
}

func (s *RecurringScheduler) Start() {
	c := cron.New()

	_, err := c.AddFunc("@every 1m", func() {
		s.processRecurring()
	})
	if err != nil {
		log.Println("scheduler registration failed:", err)
		return
	}

	c.Start()

	log.Println("recurring scheduler started")
}

func (s *RecurringScheduler) processRecurring() {
	now := time.Now()

	items, err := s.repo.FindDueRecurringTransactions(now)
	if err != nil {
		log.Println("scheduler fetch failed:", err)
		return
	}

	for _, item := range items {
		exists, err := s.repo.TransactionExistsForDate(
			item.UserID,
			item.Title,
			item.NextRunAt,
		)
		if err != nil {
			continue
		}

		if exists {
			_ = s.repo.AdvanceNextRun(&item)
			continue
		}

		transaction := model.Transaction{
			UserID:          item.UserID,
			Title:           item.Title,
			Amount:          item.Amount,
			Type:            item.Type,
			CategoryID:      item.CategoryID,
			Notes:           item.Notes,
			TransactionDate: item.NextRunAt,
		}

		if err := config.DB.Create(&transaction).Error; err != nil {
			log.Println("transaction creation failed:", err)
			continue
		}

		if err := s.repo.AdvanceNextRun(&item); err != nil {
			log.Println("advance failed:", err)
		}
	}
}
