package http

import (
	"net/http"

	"github.com/rs/cors"
)

func CorsSettings(frontHost string) *cors.Cors {
	c := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPost,
			http.MethodPut,
		},
		AllowedOrigins: []string{
			frontHost,
			"http://localhost:63342",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Accept",
			"Accept-Language",
		},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		ExposedHeaders: []string{
			"Content-Type",
			"Accept",
			"Accept-Language",
		},
		Debug: true,
	})
	return c
}
