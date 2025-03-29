package main

import (
	"github.com/JerryLegend254/mfit_api/internal/db"
	"github.com/JerryLegend254/mfit_api/internal/env"
	"github.com/JerryLegend254/mfit_api/internal/logger"
	"github.com/JerryLegend254/mfit_api/internal/store"
)

// TODO: make use of version after adding changelog for sem ver
//var version = "0.0.0"

//	@title			MFit API
//	@description	This is an API for a fitness application
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:           env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/mfit?sslmode=disable"),
			maxOpenConns:   env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns:   env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTimeout: env.GetString("DB_MAX_IDLE_TIMEOUT", "15m"),
		},
	}
	logger := logger.NewLogger()

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTimeout)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Database connection successful")

	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
