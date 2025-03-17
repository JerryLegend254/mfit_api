package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/JerryLegend254/mfit_api/internal/store"
	"github.com/go-chi/chi/v5"
)

type equipmentContextKey string

var equipmentCtxKey equipmentContextKey = "equipment"

type CreateEquipmentPayload struct {
	Name string `json:"name" validate:"required,max=40"`
}

// CreateEquipment godoc
//
//	@Summary		Creates a equipment
//	@Description	Creates a equipment
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateEquipmentPayload	true	"Equipment payload"
//	@Success		201		{object}	store.Equipment
//	@Failure		400		{object}	error
//	@Failure		403		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/equipment [post]
func (app *application) createEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload CreateEquipmentPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	equipment := store.Equipment{
		Name: payload.Name,
	}

	if err = app.store.Equipment.Create(ctx, &equipment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, &equipment); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAllEquipments godoc
//
//	@Summary		Fetch all equipment
//	@Description	Fetch all equipment
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.Equipment
//	@Failure		403	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/equipment [get]
func (app *application) fetchEquipmentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	equipment, err := app.store.Equipment.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, equipment); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetEquipment godoc
//
//	@Summary		Fetches a equipment
//	@Description	Fetches a equipment by ID
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			equipmentId	path		int	true	"Equipment ID"
//	@Success		200			{object}	store.Equipment
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/equipment/{equipmentId} [get]
func (app *application) getEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	equipment := getEquipmentFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, equipment); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeleteEquipment godoc
//
//	@Summary		Deletes a equipment
//	@Description	Deletes a equipment by ID
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Equipment ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/equipment/{equipmentId} [delete]
func (app *application) deleteEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	equipment := getEquipmentFromContext(r)

	if err := app.store.Equipment.Delete(ctx, equipment.ID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFound(w, r)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateEquipmentPayload struct {
	Name *string `json:"name" validate:"omitempty,max=40"`
}

// UddateEquipment godoc
//
//	@Summary		Update a equipment
//	@Description	Update a equipment by ID
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			equipmentId	path		int						true	"Equipment ID"
//	@Param			equipmentId	body		UpdateEquipmentPayload	true	"Equipment ID"
//	@Success		200			{object}	store.Equipment
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/equipment/{equipmentId} [patch]
func (app *application) updateEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	equipment := getEquipmentFromContext(r)

	var payload UpdateEquipmentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Name != nil {
		equipment.Name = *payload.Name
	}

	if err := app.store.Equipment.Update(ctx, equipment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, equipment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) equipmentContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParam(r, "equipmentId")

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			app.badRequest(w, r, errors.New("invalid equipment id"))
			return
		}

		equipment, err := app.store.Equipment.GetByID(ctx, intId)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFound(w, r)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, equipmentCtxKey, equipment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getEquipmentFromContext(r *http.Request) *store.Equipment {
	equipment, _ := r.Context().Value(equipmentCtxKey).(*store.Equipment)
	return equipment
}
