package main

import (
    "log"
    "net/http"

	"seriesapp/src/app/handlers"
	"seriesapp/src/database"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    db, err := database.SetupDatabase("series.db")
    if err != nil {
        log.Fatal("CRITICAL: No se pudo conectar a la base de datos:", err)
    }
    defer db.Close()

    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(corsConfiguration())
	
	//Se establecen los endpoint
	r.Route("/api", func(r chi.Router) {
		r.Route("/series", func(r chi.Router) {
			r.Get("/", handlers.GetSeriesHandler(db))  
			r.Post("/", handlers.PostSeriesHandler(db))
		})
		r.Route("/serie/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetSeriesByIdHandler(db))
			r.Put("/",handlers.PutSeriesHandler(db))
			r.Delete("/", handlers.DeleteSeriesHandler(db))
		})
	})

    port := ":8080"
    log.Printf("Servidor escuchando en puerto %s", port)
    log.Fatal(http.ListenAndServe(port, r))
}