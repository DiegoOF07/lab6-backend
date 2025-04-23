package main

import (
    "net/http"

    "github.com/go-chi/cors"
)

func corsConfiguration() func(http.Handler) http.Handler {
    middleware := cors.New(cors.Options{
        AllowedOrigins:   []string{"*", "http://localhost:5500", "http://127.0.0.1:5500"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300, 
    })
    return middleware.Handler
}
