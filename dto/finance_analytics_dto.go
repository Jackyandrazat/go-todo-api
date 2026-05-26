package dto

type CategoryAnalytics struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

type DailyCashflow struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type FinanceAnalyticsResponse struct {
	IncomeByCategory  []CategoryAnalytics `json:"income_by_category"`
	ExpenseByCategory []CategoryAnalytics `json:"expense_by_category"`
	DailyCashflow     []DailyCashflow     `json:"daily_cashflow"`
}
