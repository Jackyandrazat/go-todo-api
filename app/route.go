package app

import (
	"go-todo-api/config"
	"go-todo-api/handler"
	"go-todo-api/middleware"
	"go-todo-api/response"

	"github.com/gin-gonic/gin"

	_ "go-todo-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine) *gin.Engine {

	authHandler := handler.NewAuthHandler()
	todoHandler := handler.NewTodoHandler()
	noteHandler := handler.NewNoteHandler()
	dashboardHandler := handler.NewDashboardHandler()
	profileHandler := handler.NewProfileHandler()
	transactionHandler := handler.NewTransactionHandler()
	categoryHandler := handler.NewTransactionCategoryHandler()
	financeAnalyticsHandler := handler.NewFinanceAnalyticsHandler()
	budgetHandler := handler.NewBudgetHandler()
	recurringHandler := handler.NewRecurringTransactionHandler()
	alertHandler := handler.NewAlertHandler()

	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ready", func(c *gin.Context) {
		sqlDB, err := config.DB.DB()
		if err != nil {
			response.Error(c, 500, "db unavailable", nil)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			response.Error(c, 500, "db unavailable", nil)
			return
		}

		response.Success(c, 200, "ready", nil)
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)

		authProtected := auth.Group("")
		authProtected.Use(middleware.AuthMiddleware())
		{
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/me", authHandler.Me)
		}
	}

	todos := r.Group("/todos")
	todos.Use(middleware.AuthMiddleware())
	{
		todos.GET("", todoHandler.GetTodos)
		todos.POST("", todoHandler.CreateTodo)
		todos.PATCH("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}

	notes := r.Group("/notes")
	notes.Use(middleware.AuthMiddleware())
	{
		notes.GET("", noteHandler.GetNotes)
		notes.POST("", noteHandler.CreateNote)
		notes.PATCH("/:id", noteHandler.UpdateNote)
		notes.PATCH("/:id/pin", noteHandler.TogglePin)
		notes.DELETE("/:id", noteHandler.DeleteNote)
	}

	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware())
	{
		dashboard.GET("/summary", dashboardHandler.GetSummary)
	}

	profile := r.Group("/profile")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.GET("", profileHandler.GetProfile)
		profile.PATCH("", profileHandler.UpdateProfile)
		profile.PATCH("/password", profileHandler.ChangePassword)
	}

	transactions := r.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware())
	{
		transactions.GET("", transactionHandler.GetTransactions)
		transactions.POST("", transactionHandler.CreateTransaction)
		transactions.PATCH("/:id", transactionHandler.UpdateTransaction)
		transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
	}

	finance := r.Group("/finance")
	finance.Use(middleware.AuthMiddleware())
	{
		finance.GET("/summary", transactionHandler.GetFinanceSummary)
		finance.GET("/analytics", financeAnalyticsHandler.GetAnalytics)
	}

	categories := r.Group("/transaction-categories")
	categories.Use(middleware.AuthMiddleware())
	{
		categories.GET("", categoryHandler.GetCategories)
		categories.POST("", categoryHandler.CreateCategory)
		categories.PATCH("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	budgets := r.Group("/budgets")
	budgets.Use(middleware.AuthMiddleware())
	{
		budgets.GET("", budgetHandler.GetBudgets)
		budgets.POST("", budgetHandler.CreateBudget)
		budgets.PATCH("/:id", budgetHandler.UpdateBudget)
		budgets.DELETE("/:id", budgetHandler.DeleteBudget)
	}

	recurring := r.Group("/recurring-transactions")
	recurring.Use(middleware.AuthMiddleware())
	{
		recurring.GET("", recurringHandler.GetRecurringTransactions)
		recurring.POST("", recurringHandler.CreateRecurringTransaction)
		recurring.PATCH("/:id", recurringHandler.UpdateRecurringTransaction)
		recurring.DELETE("/:id", recurringHandler.DeleteRecurringTransaction)
	}

	alerts := r.Group("/alerts")
	alerts.Use(middleware.AuthMiddleware())
	{
		alerts.GET("", alertHandler.GetAlerts)
		alerts.PATCH("/:id/read", alertHandler.MarkAsRead)
		alerts.DELETE("/:id", alertHandler.DeleteAlert)
	}

	return r
}
