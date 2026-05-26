package service

import (
	"errors"
	"go-todo-api/dto"
	"go-todo-api/repository"
	"time"
)

type FinanceAnalyticsService struct {
	repo *repository.FinanceAnalyticsRepository
}

func NewFinanceAnalyticsService() *FinanceAnalyticsService {
	return &FinanceAnalyticsService{
		repo: repository.NewFinanceAnalyticsRepository(),
	}
}

func (s *FinanceAnalyticsService) GetAnalytics(
	userID uint,
	month string,
) (*dto.FinanceAnalyticsResponse, error) {
	parsed, err := time.Parse("2006-01", month)
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

	income, err := s.repo.GetCategoryBreakdown(
		userID,
		"income",
		start,
		end,
	)
	if err != nil {
		return nil, err
	}

	expense, err := s.repo.GetCategoryBreakdown(
		userID,
		"expense",
		start,
		end,
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.repo.GetDailyTransactions(
		userID,
		start,
		end,
	)
	if err != nil {
		return nil, err
	}

	dailyMap := map[string]*dto.DailyCashflow{}

	for _, row := range rows {
		date := row.Date.Format("2006-01-02")

		if _, exists := dailyMap[date]; !exists {
			dailyMap[date] = &dto.DailyCashflow{
				Date: date,
			}
		}

		if row.Type == "income" {
			dailyMap[date].Income = row.Amount
		} else {
			dailyMap[date].Expense = row.Amount
		}

		dailyMap[date].Balance =
			dailyMap[date].Income - dailyMap[date].Expense
	}

	var daily []dto.DailyCashflow
	for _, v := range dailyMap {
		daily = append(daily, *v)
	}

	return &dto.FinanceAnalyticsResponse{
		IncomeByCategory:  income,
		ExpenseByCategory: expense,
		DailyCashflow:     daily,
	}, nil
}
