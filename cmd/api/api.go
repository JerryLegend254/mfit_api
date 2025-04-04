package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JerryLegend254/mfit_api/docs"
	"github.com/JerryLegend254/mfit_api/internal/logger"
	"github.com/JerryLegend254/mfit_api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	config config
	store  store.Storage
	logger logger.Logger
}

type config struct {
	addr   string
	db     dbConfig
	apiURL string
}

type dbConfig struct {
	addr           string
	maxOpenConns   int
	maxIdleConns   int
	maxIdleTimeout string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Swagger Docs
	docURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docURL)))

	// Handlers
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", app.pingHandler)

		// body parts endpoints
		r.Route("/bodyparts", func(r chi.Router) {
			r.Post("/", app.createBodyPartHandler)
			r.Get("/", app.fetchBodyPartsHandler)

			r.Route("/{bodyPartId}", func(r chi.Router) {
				r.Use(app.bodyPartContextMiddleware)
				r.Get("/", app.getBodyPartHandler)
				r.Patch("/", app.updateBodyPartHandler)
				r.Delete("/", app.deleteBodyPartHandler)
			})
		})

		// targets endpoints
		r.Route("/targets", func(r chi.Router) {
			r.Post("/", app.createTargetHandler)
			r.Get("/", app.fetchTargetsHandler)

			r.Route("/{targetId}", func(r chi.Router) {
				r.Use(app.targetContextMiddleware)
				r.Get("/", app.getTargetHandler)
				r.Patch("/", app.updateTargetHandler)
				r.Delete("/", app.deleteTargetHandler)
			})
		})

		// equipment endpoints
		r.Route("/equipment", func(r chi.Router) {
			r.Post("/", app.createEquipmentHandler)
			r.Get("/", app.fetchEquipmentsHandler)

			r.Route("/{equipmentId}", func(r chi.Router) {
				r.Use(app.equipmentContextMiddleware)
				r.Get("/", app.getEquipmentHandler)
				r.Patch("/", app.updateEquipmentHandler)
				r.Delete("/", app.deleteEquipmentHandler)
			})
		})

		// workouts endpoints
		r.Route("/workouts", func(r chi.Router) {
			r.Post("/", app.createWorkoutHandler)
			r.Get("/", app.fetchWorkoutsHandler)
			r.Route("/{workoutId}", func(r chi.Router) {
				r.Use(app.workoutContextMiddleware)
				r.Get("/", app.getWorkoutHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/api/v1"

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("server started at ", srv.Addr)

	return srv.ListenAndServe()
}
