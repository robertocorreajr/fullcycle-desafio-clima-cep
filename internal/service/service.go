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
	ErrInvalidZip = errors.New("invalid zipcode")
	ErrNotFound   = errors.New("can not find zipcode")
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
		return nil, ErrInvalidZip
	}

	// consulta ViaCEP
	addr, err := s.CEP.Lookup(ctx, cep)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "400") {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("viacep: %w", err)
	}

	// Verifica se o endereço é válido
	if addr == nil || addr.Erro || addr.Localidade == "" || addr.UF == "" {
		return nil, ErrNotFound
	}

	// monta query para WeatherAPI
	q := fmt.Sprintf("%s,%s,Brazil", strings.TrimSpace(addr.Localidade), strings.TrimSpace(addr.UF))

	tempC, err := s.Weather.CurrentTempC(ctx, q)
	if err != nil {
		// Melhorara o tratamento de erros para localidades que não forem encontradas
		if strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "404") {
			return nil, ErrNotFound
		}
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
