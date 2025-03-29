package main

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest = errors.New("invalid payload check request body")
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println("internal server error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusInternalServerError, "something went wrong")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("not found", "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println("conflict error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println("bad request", "error", err.Error(), "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusBadRequest, ErrBadRequest.Error())
}
