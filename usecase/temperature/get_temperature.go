package temperature

import (
	"context"
	"encoding/json"
	"io"
	"lab-google-cloud-run/internal/dto"
	"lab-google-cloud-run/internal/entity"
	"math"
	"net/http"
	"net/url"
	"os"
	"golang.org/x/text/unicode/norm"
	"strings"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type WeatherAPIClient struct {
	httpClient HTTPDoer
	apiKey     string
	baseURL    string
}

func NewWeatherAPIClient(c HTTPDoer, apiKey string) *WeatherAPIClient {
	if c == nil {
		c = http.DefaultClient
	}
	if apiKey == "" {
		apiKey = os.Getenv("WEATHER_API_KEY")
	}
	return &WeatherAPIClient{httpClient: c, apiKey: apiKey, baseURL: "https://api.weatherapi.com/v1"}
}

func roundFloat(val float64) float64 {
	ratio := math.Pow(10, 1)
	return math.Round(val*ratio) / ratio
}

func (w *WeatherAPIClient) GetTemperatureByLocation(ctx context.Context, location string) (*entity.Response, error) {
	vals := url.Values{}
	vals.Set("q", norm.NFC.String(location))
	vals.Set("lang", "pt")
	vals.Set("key", strings.Trim(w.apiKey, "\""))
	finalURL := w.baseURL + "/current.json?" + vals.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, finalURL, nil)

	if err != nil {
		return nil, err
	}
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var weather dto.WeatherAPI
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	return &entity.Response{
		TempC: roundFloat(weather.Current.TempC),
		TempF: roundFloat(weather.Current.TempC*1.8) + 32,
		TempK: roundFloat(weather.Current.TempC + 273),
	}, nil
}

func GetTemperatureByLocation(location string) (*entity.Response, error) {
	return NewWeatherAPIClient(http.DefaultClient, "").GetTemperatureByLocation(context.Background(), location)
}
