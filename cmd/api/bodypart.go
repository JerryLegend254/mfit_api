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

type bodyPartContextKey string

var bodyPartCtxKey bodyPartContextKey = "bodyPart"

type CreateBodyPartPayload struct {
	Name     string `json:"name" validate:"required,max=40"`
	ImageUrl string `json:"image_url" validate:"required,max=255"`
}

// CreateBodyPart godoc
//
//	@Summary		Creates a body part
//	@Description	Creates a body part
//	@Tags			body parts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateBodyPartPayload	true	"BodyPart payload"
//	@Success		201		{object}	store.BodyPart
//	@Failure		400		{object}	error
//	@Failure		403		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bodyparts [post]
func (app *application) createBodyPartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload CreateBodyPartPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	bodyPart := store.BodyPart{
		Name:     payload.Name,
		ImageUrl: payload.ImageUrl,
	}

	if err = app.store.BodyParts.Create(ctx, &bodyPart); err != nil {
		switch err {
		case store.ErrDuplicate:
			app.conflictError(w, r, store.ErrDuplicateName)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, &bodyPart); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAllBodyParts godoc
//
//	@Summary		Fetch all body parts
//	@Description	Fetch all body parts
//	@Tags			body parts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.BodyPart
//	@Failure		403	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bodyparts [get]
func (app *application) fetchBodyPartsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bodyParts, err := app.store.BodyParts.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, bodyParts); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetBodyPart godoc
//
//	@Summary		Fetches a body part
//	@Description	Fetches a body part by ID
//	@Tags			body parts
//	@Accept			json
//	@Produce		json
//	@Param			bodyPartId	path		int	true	"Body Part ID"
//	@Success		200			{object}	store.BodyPart
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bodyparts/{bodyPartId} [get]
func (app *application) getBodyPartHandler(w http.ResponseWriter, r *http.Request) {
	bodyPart := getBodyPartFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, bodyPart); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeleteBodyPart godoc
//
//	@Summary		Deletes a body part
//	@Description	Deletes a body part by ID
//	@Tags			body parts
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Body Part ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bodyparts/{bodyPartId} [delete]
func (app *application) deleteBodyPartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bodyPart := getBodyPartFromContext(r)

	if err := app.store.BodyParts.Delete(ctx, bodyPart.ID); err != nil {
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

type UpdateBodyPartPayload struct {
	Name     *string `json:"name" validate:"omitempty,max=40"`
	ImageUrl *string `json:"image_url" validate:"omitempty,max=255"`
}

// UddateBodyPart godoc
//
//	@Summary		Update a body part
//	@Description	Update a body part by ID
//	@Tags			body parts
//	@Accept			json
//	@Produce		json
//	@Param			bodyPartId	path		int						true	"Body Part ID"
//	@Param			bodyPartId	body		UpdateBodyPartPayload	true	"Body Part ID"
//	@Success		200			{object}	store.BodyPart
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bodyparts/{bodyPartId} [patch]
func (app *application) updateBodyPartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bodyPart := getBodyPartFromContext(r)

	var payload UpdateBodyPartPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Name != nil {
		bodyPart.Name = *payload.Name
	}

	if payload.ImageUrl != nil {
		bodyPart.ImageUrl = *payload.ImageUrl
	}

	if err := app.store.BodyParts.Update(ctx, bodyPart); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, bodyPart); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) bodyPartContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParam(r, "bodyPartId")

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			app.badRequest(w, r, errors.New("invalid body part id"))
			return
		}

		bodyPart, err := app.store.BodyParts.GetByID(ctx, intId)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFound(w, r)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, bodyPartCtxKey, bodyPart)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getBodyPartFromContext(r *http.Request) *store.BodyPart {
	bodyPart, _ := r.Context().Value(bodyPartCtxKey).(*store.BodyPart)
	return bodyPart
}
