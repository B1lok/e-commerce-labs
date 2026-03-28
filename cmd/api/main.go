package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"weather-api/internal/application/scheduled"
	"weather-api/internal/application/services"
	"weather-api/internal/config"
	postgresconnector "weather-api/internal/infrastructure/db/postgres"
	"weather-api/internal/infrastructure/email"
	cityValidator "weather-api/internal/infrastructure/http/validator"
	weatherapi "weather-api/internal/infrastructure/http/weather-api"
	"weather-api/internal/interface/api/rest"
	"weather-api/pkg/logger"
	"weather-api/pkg/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	postgresconnector.RunMigrations(cfg)

	db, err := postgresconnector.ConnectDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	weatherRepo := weatherapi.NewWeatherRepository(cfg.WeatherApiKey)
	weatherService := services.NewWeatherService(weatherRepo)
	weatherController := rest.NewWeatherController(weatherService)
	cityValidatorImpl := cityValidator.NewCityValidator(cfg.WeatherApiKey)
	sender := email.NewEmailSender(email.CreateConfig(cfg))
	txManager := middleware.NewTxManager(db)

	subscriptionRepo := postgresconnector.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, cityValidatorImpl, sender, cfg.ServerHost)
	subscriptionController := rest.NewSubscriptionController(subscriptionService)

	jm := scheduled.NewJobManager(ctx)
	jm.RegisterJob(scheduled.NewHourlyWeatherUpdateJob(weatherRepo, subscriptionRepo, sender, cfg.ServerHost))
	jm.RegisterJob(scheduled.NewDailyWeatherUpdateJob(weatherRepo, subscriptionRepo, sender, cfg.ServerHost))
	go jm.StartScheduler()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("templates/index.html")
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.TransactionMiddleware(txManager))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			slog.Error("Health check failed: cannot get underlying DB", "error", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "reason": "database connection error"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			slog.Error("Health check failed: database ping failed", "error", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "reason": "database unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := router.Group("/api")
	{
		api.GET("/weather", weatherController.GetWeather)
		api.POST("/subscribe", subscriptionController.Subscribe)
		api.GET("/confirm/:token", subscriptionController.Confirm)
		api.GET("/unsubscribe/:token", subscriptionController.Unsubscribe)
	}

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		slog.Info("Server starting", "address", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("SIGTERM received. Starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	} else {
		slog.Info("HTTP server stopped gracefully")
	}

	jm.Stop()

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("Error closing database connection", "error", err)
		} else {
			slog.Info("Database connection closed")
		}
	}

	slog.Info("Shutdown complete. Exiting.")
}
