package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/JerryLegend254/mfit_api/internal/logger"
	"github.com/JerryLegend254/mfit_api/internal/store"
)

func newTestApplication(t testing.TB, store store.Storage) *application {
	t.Helper()

	logger := logger.NewLogger()

	return &application{
		store:  store,
		logger: logger,
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d want %d", got, want)

	}
}

func assertResponse(t testing.TB, actualBody *bytes.Buffer, expectedBody []byte) {
	t.Helper()

	var expectedMap, actualMap map[string]interface{}

	// Convert actual response to a map
	json.Unmarshal(actualBody.Bytes(), &actualMap)

	// Convert expected response to a map (if it's a string, first convert it)
	json.Unmarshal(expectedBody, &expectedMap)

	if !reflect.DeepEqual(expectedMap, actualMap) {
		t.Errorf("got %v want %v", actualMap, expectedBody)

	}
}

func assertContentType(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)

	}
}
func newTestDB(t testing.TB) (*sql.DB, func()) {
	db, err := sql.Open("postgres", "postgres://admin:adminpassword@localhost/mfit_test?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to start test db: %v", err)
	}

	m, err := runMigrations(db, "file://../migrate/migrations")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	teardown := func() {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			t.Errorf("Failed to roll back migrations: %v", err)
		}
		db.Close()
	}

	return db, teardown
}

func runMigrations(db *sql.DB, migrationsPath string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return m, nil
}
