package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"os/signal"
	"syscall"
	"weather-api/internal/application/scheduled"
	"weather-api/internal/application/services"
	"weather-api/internal/config"
	postgresconnector "weather-api/internal/infrastructure/db/postgres"
	"weather-api/internal/infrastructure/email"
	cityValidator "weather-api/internal/infrastructure/http/validator"
	weatherapi "weather-api/internal/infrastructure/http/weather-api"
	"weather-api/internal/interface/api/rest"
	"weather-api/pkg/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	postgresconnector.RunMigrations(cfg)

	db, err := postgresconnector.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

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

	router := gin.Default()
	router.LoadHTMLGlob("templates/index.html")
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.TransactionMiddleware(txManager))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	api := router.Group("/api")
	{
		api.GET("/weather", weatherController.GetWeather)
		api.POST("/subscribe", subscriptionController.Subscribe)
		api.GET("/confirm/:token", subscriptionController.Confirm)
		api.GET("/unsubscribe/:token", subscriptionController.Unsubscribe)
	}

	go func() {
		<-ctx.Done()
		cancel()
	}()

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
