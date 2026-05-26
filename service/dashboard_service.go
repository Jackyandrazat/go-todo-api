package service

import (
	"errors"
	"go-todo-api/repository"
	"time"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

type DashboardSummary struct {
	Tasks       *repository.TaskStats      `json:"tasks"`
	Notes       *repository.NoteStats      `json:"notes"`
	Finance     *repository.FinanceSummary `json:"finance"`
	RecentTasks interface{}                `json:"recent_tasks"`
	RecentNotes interface{}                `json:"recent_notes"`
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		repo: repository.NewDashboardRepository(),
	}
}

func (s *DashboardService) GetSummary(
	userID uint,
	month *string,
) (*DashboardSummary, error) {
	taskStats, err := s.repo.GetTaskStats(userID)
	if err != nil {
		return nil, err
	}

	noteStats, err := s.repo.GetNoteStats(userID)
	if err != nil {
		return nil, err
	}

	recentTasks, err := s.repo.GetRecentTasks(userID, 5)
	if err != nil {
		return nil, err
	}

	recentNotes, err := s.repo.GetRecentNotes(userID, 5)
	if err != nil {
		return nil, err
	}

	var startDate *time.Time
	var endDate *time.Time

	if month != nil {
		parsed, err := time.Parse("2006-01", *month)
		if err != nil {
			return nil, errors.New("invalid month format, use YYYY-MM")
		}

		start := time.Date(
			parsed.Year(),
			parsed.Month(),
			1,
			0, 0, 0, 0,
			time.UTC,
		)

		end := start.AddDate(0, 1, -1)

		startDate = &start
		endDate = &end
	}

	finance, err := s.repo.GetFinanceSummary(
		userID,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}

	return &DashboardSummary{
		Tasks:       taskStats,
		Notes:       noteStats,
		Finance:     finance,
		RecentTasks: recentTasks,
		RecentNotes: recentNotes,
	}, nil
}
