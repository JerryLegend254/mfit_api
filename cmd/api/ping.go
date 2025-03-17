package main

import "net/http"

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"message": "Welcome to MFit",
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
