package weather_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"weather-api/internal/domain/models"
	"weather-api/pkg/errors"
)

type WeatherRepository struct {
	apiKey string
	client *http.Client
	clock  Clock
}

func NewWeatherRepository(apiKey string) *WeatherRepository {
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return &WeatherRepository{apiKey: apiKey, client: client, clock: SystemClock{}}
}

const (
	baseUrl          = "http://api.weatherapi.com/v1"
	currentEndpoint  = "/current.json"
	forecastEndpoint = "/forecast.json"
	defaultTimeout   = 10 * time.Second
)

func (r *WeatherRepository) SetClock(clock Clock) {
	r.clock = clock
}

func (r *WeatherRepository) GetWeather(ctx context.Context, city string) (*models.Weather, error) {
	endpoint := fmt.Sprintf("%s%s?key=%s&q=%s",
		baseUrl,
		currentEndpoint,
		r.apiKey,
		url.QueryEscape(city),
	)
	resp, err := r.requestWeatherAPI(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := r.handleAPIResponse(resp); err != nil {
		return nil, err
	}
	var apiResponse WeatherRepositoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse weather data", http.StatusInternalServerError)
	}
	return ToWeather(&apiResponse), nil
}

func (r *WeatherRepository) GetDailyForecast(ctx context.Context, city string) (*models.WeatherDaily, error) {
	endpoint := fmt.Sprintf("%s%s?key=%s&q=%s&days=1",
		baseUrl,
		forecastEndpoint,
		r.apiKey,
		url.QueryEscape(city),
	)
	resp, err := r.requestWeatherAPI(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := r.handleAPIResponse(resp); err != nil {
		return nil, err
	}

	var apiResponse WeatherDailyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse weather data", http.StatusInternalServerError)
	}
	return ToWeatherDaily(&apiResponse), nil
}

func (r *WeatherRepository) GetHourlyForecast(ctx context.Context, city string) (*models.WeatherHourly, error) {
	endpoint := fmt.Sprintf("%s%s?key=%s&q=%s&days=1",
		baseUrl,
		forecastEndpoint,
		r.apiKey,
		url.QueryEscape(city),
	)
	resp, err := r.requestWeatherAPI(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := r.handleAPIResponse(resp); err != nil {
		return nil, err
	}

	var apiResponse WeatherHourlyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse weather data", http.StatusInternalServerError)
	}
	return ToWeatherHourly(&apiResponse, r.clock.Now()), nil
}
func (r *WeatherRepository) requestWeatherAPI(ctx context.Context, endpoint string) (*http.Response, error) {
	resp, err := r.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to weather API", http.StatusServiceUnavailable)
	}
	return resp, nil
}

func (r *WeatherRepository) handleAPIResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return errors.Wrap(err, "failed to parse error response from weather API", http.StatusBadGateway)
		}

		if errResp.Error.Code == 1006 {
			return errors.New("City not found", http.StatusNotFound)
		}

		return errors.New(fmt.Sprintf("weather API error: %s", errResp.Error.Message), http.StatusBadGateway)
	}

	return nil
}
