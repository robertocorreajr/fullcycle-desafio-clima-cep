package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/http"
	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/service"
	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/types"
)

// --- mocks ---

type mockCEP struct {
	localidade string
	uf         string
	err        error
	notFound   bool
}

func (m mockCEP) Lookup(ctx context.Context, cep string) (*types.ViaCEPResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.notFound {
		return &types.ViaCEPResult{Erro: true}, nil
	}
	return &types.ViaCEPResult{Localidade: m.localidade, UF: m.uf}, nil
}

type mockWeather struct {
	tempC float64
	err   error
}

func (m mockWeather) CurrentTempC(ctx context.Context, query string) (float64, error) {
	return m.tempC, m.err
}

// --- tests ---

func TestSuccess(t *testing.T) {
	svc := service.New(mockCEP{localidade: "São Paulo", uf: "SP"}, mockWeather{tempC: 28.5})
	h := apphttp.NewRouter(&apphttp.Handler{Svc: svc})

	req := httptest.NewRequest(http.MethodGet, "/weather/01001000", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", rr.Code)
	}

	var body types.WeatherResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if body.TempC != 28.5 {
		t.Fatalf("unexpected tempC: %v", body.TempC)
	}
}

func TestInvalidZip(t *testing.T) {
	svc := service.New(mockCEP{}, mockWeather{})
	h := apphttp.NewRouter(&apphttp.Handler{Svc: svc})

	req := httptest.NewRequest(http.MethodGet, "/weather/123", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want 422 got %d", rr.Code)
	}
}

func TestNotFound(t *testing.T) {
	svc := service.New(mockCEP{notFound: true}, mockWeather{})
	h := apphttp.NewRouter(&apphttp.Handler{Svc: svc})

	req := httptest.NewRequest(http.MethodGet, "/weather/99999999", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", rr.Code)
	}
}
