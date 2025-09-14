package main

import (
	"lab-google-cloud-run/infra/webserver/handlers"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	_, ok := os.LookupEnv("WEATHER_API_KEY")
	if !ok {
		panic("WEATHER_API_KEY not set")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/cep-temperature", func(r chi.Router) {
		r.Method(http.MethodGet, "/{cep}", handlers.GetCepTemperatureHandler())
	})
	http.ListenAndServe(":8080", r)
}
