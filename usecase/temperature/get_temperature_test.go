package temperature

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

func newTestClient(fn roundTripFunc) HTTPDoer { return roundTripFunc(fn) }

func TestWeatherAPIClient_Success(t *testing.T) {
	bodyJSON := `{"current":{"temp_c":25.0,"temp_f":77.0}}`
	wClient := NewWeatherAPIClient(newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(bodyJSON))}, nil
	}), "dummy")
	resp, err := wClient.GetTemperatureByLocation(context.Background(), "Sao Paulo")
	if err != nil {
		t.Fatalf("erro: %v", err)
	}
	if resp.TempC != 25.0 {
		t.Errorf("TempC esperada 25.0 got %v", resp.TempC)
	}
}

func TestWeatherAPIClient_HTTPError(t *testing.T) {
	wClient := NewWeatherAPIClient(newTestClient(func(r *http.Request) (*http.Response, error) { return nil, errors.New("http fail") }), "dummy")
	_, err := wClient.GetTemperatureByLocation(context.Background(), "Sao Paulo")
	if err == nil {
		t.Fatalf("esperava erro")
	}
}

func TestWeatherAPIClient_InvalidJSON(t *testing.T) {
	wClient := NewWeatherAPIClient(newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader("{\"current\":"))}, nil
	}), "dummy")
	_, err := wClient.GetTemperatureByLocation(context.Background(), "Sao Paulo")
	if err == nil {
		t.Fatalf("esperava erro de JSON")
	}
}
