package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"lab-google-cloud-run/internal/entity"
	"lab-google-cloud-run/pkg"
	"lab-google-cloud-run/usecase/buscacep"
	"lab-google-cloud-run/usecase/temperature"
	"net/http"

	"github.com/go-chi/chi"
)

type CepLocator interface {
	GetLocationByCep(ctx context.Context, cep string) (string, error)
}

type TemperatureProvider interface {
	GetTemperatureByLocation(ctx context.Context, location string) (*entity.Response, error)
}

type DefaultCepLocator struct{}

func (d DefaultCepLocator) GetLocationByCep(ctx context.Context, cep string) (string, error) {
	// Usa client com contexto
	client := buscacep.NewViaCepClient(http.DefaultClient)
	return client.GetLocationByCep(ctx, cep)
}

type DefaultTemperatureProvider struct{}

func (d DefaultTemperatureProvider) GetTemperatureByLocation(ctx context.Context, location string) (*entity.Response, error) {
	client := temperature.NewWeatherAPIClient(http.DefaultClient, "")
	return client.GetTemperatureByLocation(ctx, location)
}

type CepTemperatureHandler struct {
	locator     CepLocator
	temperature TemperatureProvider
}

func NewCepTemperatureHandler(locator CepLocator, temperature TemperatureProvider) *CepTemperatureHandler {
	return &CepTemperatureHandler{locator: locator, temperature: temperature}
}

func (h *CepTemperatureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cepParam := chi.URLParam(r, "cep")
	ctx := r.Context()

	if !pkg.ValidateCepFormat(cepParam) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&entity.ErrorResponse{Message: "invalid zip code"})
		return
	}

	location, err := h.locator.GetLocationByCep(ctx, cepParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&entity.ErrorResponse{Message: fmt.Sprintf("internal error: %s", err.Error())})
		return
	}
	if location == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&entity.ErrorResponse{Message: "can not find zipcode"})
		return
	}

	resp, err := h.temperature.GetTemperatureByLocation(ctx, location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&entity.ErrorResponse{Message: fmt.Sprintf("internal error: %s", err.Error())})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetCepTemperatureHandler() http.Handler {
	return NewCepTemperatureHandler(DefaultCepLocator{}, DefaultTemperatureProvider{})
}
