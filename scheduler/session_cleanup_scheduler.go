package scheduler

import (
	"log"

	"go-todo-api/repository"

	"github.com/robfig/cron/v3"
)

type SessionCleanupScheduler struct {
	repo *repository.SessionRepository
}

func NewSessionCleanupScheduler() *SessionCleanupScheduler {
	return &SessionCleanupScheduler{
		repo: repository.NewSessionRepository(),
	}
}

func (s *SessionCleanupScheduler) Start() {
	c := cron.New()

	_, err := c.AddFunc("@daily", func() {
		if err := s.repo.CleanupExpired(); err != nil {
			log.Println(err)
		}
	})

	if err != nil {
		return
	}

	c.Start()
}
