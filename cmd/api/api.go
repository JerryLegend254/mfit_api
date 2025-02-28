package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct{}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Handlers
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", app.pingHandler)

		// workouts endpoints
		r.Route("/workouts", func(r chi.Router) {
			r.Post("/", app.createWorkoutHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server started at %s", srv.Addr)

	return srv.ListenAndServe()
}
