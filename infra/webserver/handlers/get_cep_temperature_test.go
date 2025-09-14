package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"lab-google-cloud-run/internal/entity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

// Mocks
type mockLocator struct {
	loc string
	err error
}

func (m mockLocator) GetLocationByCep(ctx context.Context, cep string) (string, error) {
	return m.loc, m.err
}

type mockTemp struct {
	resp *entity.Response
	err  error
}

func (m mockTemp) GetTemperatureByLocation(ctx context.Context, location string) (*entity.Response, error) {
	return m.resp, m.err
}

// helper
func performRequest(h http.Handler, cep string) *httptest.ResponseRecorder {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/cep/{cep}/temperature", h)
	req := httptest.NewRequest(http.MethodGet, "/cep/"+cep+"/temperature", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("cep", cep)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestHandler_Success(t *testing.T) {
	h := NewCepTemperatureHandler(mockLocator{loc: "Sao Paulo"}, mockTemp{resp: &entity.Response{TempC: 21.1, TempF: 70.0, TempK: 294.1}})
	rr := performRequest(h, "01001000")
	if rr.Code != http.StatusOK {
		t.Fatalf("esperado 200 got %d", rr.Code)
	}
	var resp entity.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.TempC != 21.1 {
		t.Errorf("TempC esperada 21.1 got %v", resp.TempC)
	}
}

func TestHandler_InvalidCep(t *testing.T) {
	h := NewCepTemperatureHandler(mockLocator{}, mockTemp{})
	rr := performRequest(h, "123")
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("esperado 422 got %d", rr.Code)
	}
}

func TestHandler_CepNotFound(t *testing.T) {
	h := NewCepTemperatureHandler(mockLocator{loc: ""}, mockTemp{})
	rr := performRequest(h, "01001000")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("esperado 404 got %d", rr.Code)
	}
}

func TestHandler_ErrorCepLookup(t *testing.T) {
	h := NewCepTemperatureHandler(mockLocator{err: errors.New("fail")}, mockTemp{})
	rr := performRequest(h, "01001000")
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("esperado 500 got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "internal error") {
		t.Errorf("mensagem não contém internal error: %s", rr.Body.String())
	}
}

func TestHandler_ErrorTemperature(t *testing.T) {
	h := NewCepTemperatureHandler(mockLocator{loc: "Sao Paulo"}, mockTemp{err: errors.New("temp fail")})
	rr := performRequest(h, "01001000")
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("esperado 500 got %d", rr.Code)
	}
}
