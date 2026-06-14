package scheduler

import (
	"fmt"
	"log"
	"time"

	"go-todo-api/model"
	"go-todo-api/repository"

	"github.com/robfig/cron/v3"
)

type AlertScheduler struct {
	alertRepo  *repository.AlertRepository
	budgetRepo *repository.BudgetRepository
}

func NewAlertScheduler() *AlertScheduler {
	return &AlertScheduler{
		alertRepo:  repository.NewAlertRepository(),
		budgetRepo: repository.NewBudgetRepository(),
	}
}

func (s *AlertScheduler) Start() {
	c := cron.New()

	_, err := c.AddFunc("@every 30m", func() {
		s.ProcessBudgetAlerts()
	})
	if err != nil {
		log.Println(err)
		return
	}

	c.Start()
}

func (s *AlertScheduler) ProcessBudgetAlerts() {
	now := time.Now()
	currentMonth := now.Format("2006-01")

	budgets, err := s.budgetRepo.GetActiveBudgetsWithUsage(currentMonth)
	if err != nil {
		log.Println("failed to fetch active budgets with usage for alerts:", err)
		return
	}

	for _, b := range budgets {
		if b.Amount <= 0 {
			continue
		}

		ratio := b.Spent / b.Amount

		var alertType string
		var title string
		var message string

		if ratio >= 1.0 {
			alertType = "budget_exceeded"
			title = fmt.Sprintf("Anggaran %s Terlampaui!", b.CategoryName)
			message = fmt.Sprintf("Pengeluaran Anda untuk kategori %s di bulan %s telah mencapai Rp%.2f (%.1f%%) dari batas anggaran sebesar Rp%.2f.",
				b.CategoryName, b.Month, b.Spent, ratio*100, b.Amount)
		} else if ratio >= 0.8 {
			alertType = "budget_warning"
			title = fmt.Sprintf("Peringatan Anggaran %s!", b.CategoryName)
			message = fmt.Sprintf("Pengeluaran Anda untuk kategori %s di bulan %s telah mencapai Rp%.2f (%.1f%%) dari batas anggaran sebesar Rp%.2f.",
				b.CategoryName, b.Month, b.Spent, ratio*100, b.Amount)
		} else {
			continue
		}

		// Check if there is already an unread alert of the same type and title for this user
		exists, err := s.alertRepo.ExistsSameUnread(b.UserID, alertType, title)
		if err != nil {
			log.Println("failed to check existing unread alerts:", err)
			continue
		}

		if exists {
			continue
		}

		// Create the alert
		newAlert := model.Alert{
			UserID:  b.UserID,
			Type:    alertType,
			Title:   title,
			Message: message,
			IsRead:  false,
		}

		if err := s.alertRepo.Create(&newAlert); err != nil {
			log.Println("failed to create budget alert:", err)
		}
	}
}
