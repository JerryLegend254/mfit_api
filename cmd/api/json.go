package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := int64(1 << 20) // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type errorResponse struct {
		Error string `json:"error"`
	}

	data := &errorResponse{
		Error: message,
	}

	return writeJSON(w, status, data)
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	type jsonResponse struct {
		Data interface{} `json:"data"`
	}
	return writeJSON(w, status, &jsonResponse{Data: data})
}
