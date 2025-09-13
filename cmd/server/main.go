package main

import (
	"log"
	"net/http"
	"os"

	h "github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/http"
	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/service"
	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/viacep"
	"github.com/robertocorreajr/fullcycle-desafio-clima-cep/internal/weather"
)

func main() {
	apiKey := os.Getenv("WEATHERAPI_KEY")
	if apiKey == "" {
		log.Fatal("missing WEATHERAPI_KEY")
	}

	cepClient := viacep.New()
	weatherClient := weather.New(apiKey)
	svc := service.New(cepClient, weatherClient)

	router := h.NewRouter(&h.Handler{Svc: svc})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
