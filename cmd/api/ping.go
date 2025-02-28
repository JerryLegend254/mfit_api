package main

import "net/http"

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to mfit"))
}
