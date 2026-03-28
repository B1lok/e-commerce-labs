package integration

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-api/internal/application/services"
	"weather-api/internal/interface/api/rest"
	"weather-api/internal/test/stubs"
	"weather-api/pkg/middleware"
)

type WeatherControllerTestSuite struct {
	suite.Suite
	Router *gin.Engine
}

func (suite *WeatherControllerTestSuite) SetupSuite() {
	weatherRepo := stubs.NewWeatherRepositoryStub()
	weatherService := services.NewWeatherService(weatherRepo)
	weatherController := rest.NewWeatherController(weatherService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	api := router.Group("/api")
	api.GET("/weather", weatherController.GetWeather)

	suite.Router = router
}

func (suite *WeatherControllerTestSuite) TestGetWeather() {
	req, _ := http.NewRequest("GET", "/api/weather?city=London", nil)
	resp := httptest.NewRecorder()

	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Contains(response, "description")
	suite.Contains(response, "temperature")
	suite.Contains(response, "humidity")
}

func (suite *WeatherControllerTestSuite) TestGetWeatherInvalidCity() {
	req, _ := http.NewRequest("GET", "/api/weather?city=InvalidCity", nil)
	resp := httptest.NewRecorder()

	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusNotFound, resp.Code)
}

func (suite *WeatherControllerTestSuite) TestGetWeatherInvalidQueryParam() {
	req, _ := http.NewRequest("GET", "/api/weather?city=", nil)
	resp := httptest.NewRecorder()

	suite.Router.ServeHTTP(resp, req)

	suite.Equal(http.StatusBadRequest, resp.Code)
}

func TestWeatherControllerTestSuite(t *testing.T) {
	suite.Run(t, new(WeatherControllerTestSuite))
}
