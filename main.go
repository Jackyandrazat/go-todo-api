package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-todo-api/app"
	"go-todo-api/config"
	"go-todo-api/middleware"
	"go-todo-api/scheduler"
	"go-todo-api/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()

	utils.InitLogger(config.Config.AppEnv)
	defer utils.Logger.Sync()

	if config.Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	recurringScheduler := scheduler.NewRecurringScheduler()
	recurringScheduler.Start()

	alertScheduler := scheduler.NewAlertScheduler()
	alertScheduler.Start()

	cleanupScheduler := scheduler.NewSessionCleanupScheduler()
	cleanupScheduler.Start()

	r := gin.New()

	r.Use(gin.Logger())

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.StructuredLogger())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		utils.Logger.Fatal(
			"failed to set trusted proxies",
			zap.Error(err),
		)
	}

	app.SetupRouter(r)

	server := &http.Server{
		Addr:              ":" + config.Config.AppPort,
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		utils.Logger.Info(
			"server started",
			zap.String("port", config.Config.AppPort),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			utils.Logger.Fatal(
				"server failed",
				zap.Error(err),
			)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	utils.Logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		utils.Logger.Fatal(
			"forced shutdown",
			zap.Error(err),
		)
	}

	utils.Logger.Info("server exited cleanly")
}
