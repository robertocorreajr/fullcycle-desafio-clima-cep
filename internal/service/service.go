package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/types"
	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/viacep"
	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/weather"
)

var (
	errInvalidZip = errors.New("invalid zipcode")
	errNotFound   = errors.New("can not find zipcode")
	cepRe         = regexp.MustCompile(`^\d{8}$`)
)

type Service struct {
	CEP     viacep.Client
	Weather weather.Client
}

func New(cep viacep.Client, w weather.Client) *Service {
	return &Service{CEP: cep, Weather: w}
}

func (s *Service) GetWeatherByCEP(ctx context.Context, cep string) (*types.WeatherResponse, error) {
	// valida CEP (8 dígitos)
	if !cepRe.MatchString(cep) {
		return nil, errInvalidZip
	}

	// consulta ViaCEP
	addr, err := s.CEP.Lookup(ctx, cep)
	if err != nil {
		return nil, fmt.Errorf("viacep: %w", err)
	}
	if addr == nil || addr.Erro || addr.Localidade == "" || addr.UF == "" {
		return nil, errNotFound
	}

	// monta query para WeatherAPI
	q := fmt.Sprintf("%s,%s,Brazil", strings.TrimSpace(addr.Localidade), strings.TrimSpace(addr.UF))

	tempC, err := s.Weather.CurrentTempC(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("weatherapi: %w", err)
	}

	// Conversões conforme informado no desafio
	tempF := tempC*1.8 + 32
	tempK := tempC + 273 // (observação: pesquisei e o correto fisicamente é 273.15, mas vou seguir com o informado no desafio!)

	return &types.WeatherResponse{
		TempC: round1(tempC),
		TempF: round1(tempF),
		TempK: round1(tempK),
	}, nil
}

// arredonda com 1 casa p/ stabilizar resposta (ex.: 28.5)
func round1(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

// Erros públicos para o handler decidir status code
func ErrInvalidZip() error { return errInvalidZip }
func ErrNotFound() error   { return errNotFound }
