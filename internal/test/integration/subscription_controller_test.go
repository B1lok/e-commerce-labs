package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"weather-api/internal/application/services"
	"weather-api/internal/config"
	"weather-api/internal/domain/models"
	postgresconnector "weather-api/internal/infrastructure/db/postgres"
	"weather-api/internal/interface/api/rest"
	"weather-api/internal/test/containers"
	"weather-api/internal/test/stubs"
	"weather-api/pkg/middleware"
)

type SubscriptionControllerTestSuite struct {
	suite.Suite
	Router   *gin.Engine
	Postgres *containers.PostgresContainer
	DB       *gorm.DB
}

func (suite *SubscriptionControllerTestSuite) SetupSuite() {
	ctx := context.Background()

	postgres, err := containers.SetupPostgresContainer(ctx)
	suite.Require().NoError(err)
	suite.Postgres = postgres
	host, err := postgres.Container.Host(ctx)
	suite.Require().NoError(err)

	mappedPort, err := postgres.Container.MappedPort(ctx, "5432")
	suite.Require().NoError(err)

	cfg := config.Config{
		DBUser:     "test",
		DBPassword: "test",
		DBHost:     host,
		DBPort:     mappedPort.Port(),
		DBName:     "testdb",
	}

	db, err := postgresconnector.ConnectDB(cfg)
	suite.Require().NoError(err)
	suite.DB = db
	postgresconnector.RunMigrationsWithPath(cfg, getMigrationPath())

	cityValidatorImpl := stubs.NewCityValidatorStub()
	sender := stubs.NewSenderStub()
	subscriptionRepo := postgresconnector.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, cityValidatorImpl, sender, cfg.ServerHost)
	subscriptionController := rest.NewSubscriptionController(subscriptionService)
	txManager := middleware.NewTxManager(db)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.TransactionMiddleware(txManager))

	api := router.Group("/api")
	{
		api.POST("/subscribe", subscriptionController.Subscribe)
		api.GET("/confirm/:token", subscriptionController.Confirm)
		api.GET("/unsubscribe/:token", subscriptionController.Unsubscribe)
	}

	suite.Router = router
}

func (suite *SubscriptionControllerTestSuite) TearDownTest() {
	result := suite.DB.Exec("DELETE FROM subscriptions")
	suite.Require().NoError(result.Error)
}

func (suite *SubscriptionControllerTestSuite) insertTestData(data interface{}) {
	result := suite.DB.Create(data)
	suite.Require().NoError(result.Error)
}

func (suite *SubscriptionControllerTestSuite) TestSubscribe() {
	formData := "email=test@example.com&city=London&frequency=daily"
	body := strings.NewReader(formData)

	req, _ := http.NewRequest("POST", "/api/subscribe", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)

	var count int64
	err := suite.DB.Model(&models.Subscription{}).Where("email = ?", "test@example.com").Count(&count).Error
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)
}
func (suite *SubscriptionControllerTestSuite) TestSubscribe_InvalidInput() {
	formData := "email=test@example.com"
	body := strings.NewReader(formData)

	req, _ := http.NewRequest("POST", "/api/subscribe", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Contains(resp.Body.String(), "Invalid input")
}

func (suite *SubscriptionControllerTestSuite) TestSubscribe_EmailAlreadySubscribed() {
	formData := "email=test@example.com&city=London&frequency=daily"
	body := strings.NewReader(formData)

	req1, _ := http.NewRequest("POST", "/api/subscribe", body)
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp1 := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp1, req1)
	suite.Equal(http.StatusOK, resp1.Code)

	bodyDuplicate := strings.NewReader(formData)
	req2, _ := http.NewRequest("POST", "/api/subscribe", bodyDuplicate)
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp2 := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp2, req2)

	suite.Equal(http.StatusConflict, resp2.Code)
	suite.Contains(resp2.Body.String(), "Email already subscribed")
}

func (suite *SubscriptionControllerTestSuite) TestConfirmSubscription() {
	token := "test-token"
	suite.insertTestData(&models.Subscription{
		Email:     "test@example.com",
		City:      "London",
		Frequency: "daily",
		Token:     token,
		Confirmed: false,
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/confirm/%s", token), nil)
	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Contains(response, "message")

	var subscription models.Subscription
	err = suite.DB.Where("token = ?", token).First(&subscription).Error
	suite.Require().NoError(err)
	suite.True(subscription.Confirmed)
}

func (suite *SubscriptionControllerTestSuite) TestConfirmSubscription_InvalidToken() {
	req, _ := http.NewRequest("GET", "/api/confirm/ ", nil)
	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Contains(resp.Body.String(), "Invalid token")
}

func (suite *SubscriptionControllerTestSuite) TestConfirmSubscription_TokenNotFound() {
	nonExistentToken := "non-existent-token"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/confirm/%s", nonExistentToken), nil)
	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusNotFound, resp.Code)
	suite.Contains(resp.Body.String(), "Token not found")
}

func (suite *SubscriptionControllerTestSuite) TestUnsubscribe() {
	token := "test-token"
	suite.insertTestData(&models.Subscription{
		Email:     "test@example.com",
		City:      "London",
		Frequency: "daily",
		Token:     token,
		Confirmed: true,
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/unsubscribe/%s", token), nil)
	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Contains(response, "message")

	var count int64
	err = suite.DB.Model(&models.Subscription{}).Where("token = ?", token).Count(&count).Error
	suite.Require().NoError(err)
	suite.Equal(int64(0), count)
}

func (suite *SubscriptionControllerTestSuite) TestUnsubscribe_InvalidToken() {
	req, _ := http.NewRequest("GET", "/api/unsubscribe/ ", nil)
	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Contains(resp.Body.String(), "Invalid token")
}

func (suite *SubscriptionControllerTestSuite) TestUnsubscribe_TokenNotFound() {
	nonExistentToken := "non-existent-token"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/unsubscribe/%s", nonExistentToken), nil)
	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusNotFound, resp.Code)
	suite.Contains(resp.Body.String(), "Token not found")
}

func TestSubscriptionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionControllerTestSuite))
}

func getMigrationPath() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	projectRoot := filepath.Join(workingDir, "../../..")
	migrationsPath := filepath.Join(projectRoot, "migrations")
	return fmt.Sprintf("file://%s", filepath.ToSlash(migrationsPath))
}
