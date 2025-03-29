package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/JerryLegend254/mfit_api/internal/logger"
	"github.com/JerryLegend254/mfit_api/internal/store"
)

func newTestApplication(t *testing.T, store store.Storage) *application {
	t.Helper()

	logger := logger.NewLogger()

	return &application{
		store:  store,
		logger: logger,
	}
}

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d want %d", got, want)

	}
}

func assertResponse(t *testing.T, actualBody *bytes.Buffer, expectedBody []byte) {
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
