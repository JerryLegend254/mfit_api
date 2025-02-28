package main

import (
	"encoding/json"
	"net/http"
)

type CreateWorkoutPayload struct {
	Name              string `json:"name" validate:"required,max=40"`
	Category          string `json:"category" validate:"required"`
	Difficulty        string `json:"difficulty" validate:"required"`
	CaloriesBurned    int32  `json:"calories_burned" validate:"required"`
	DurationMinutes   int32  `json:"duration_minutes" validate:"required"`
	EquipmentRequired bool   `json:"equipment_required" validate:"required"`
}

func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateWorkoutPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, payload); err != nil {
		app.internalServerError(w, r, err)
	}
}
