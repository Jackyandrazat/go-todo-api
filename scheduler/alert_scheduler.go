package scheduler

import (
	"log"

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
		s.processBudgetAlerts()
	})
	if err != nil {
		log.Println(err)
		return
	}

	c.Start()
}

func (s *AlertScheduler) processBudgetAlerts() {
	panic("unimplemented")
	// MVP implementation:
	// iterate active month budgets
	// if spent >= 80% => warning
	// if spent >= 100% => exceeded
}
