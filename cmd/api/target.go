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

type targetContextKey string

var targetCtxKey targetContextKey = "target"

type CreateTargetPayload struct {
	Name       string `json:"name" validate:"required,max=40"`
	BodyPartID int64  `json:"bodypart_id" validate:"required"`
}

// CreateTarget godoc
//
//	@Summary		Creates a target
//	@Description	Creates a target
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateTargetPayload	true	"Target payload"
//	@Success		201		{object}	store.Target
//	@Failure		400		{object}	error
//	@Failure		403		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/targets [post]
func (app *application) createTargetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload CreateTargetPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	target := store.Target{
		Name:       payload.Name,
		BodyPartID: payload.BodyPartID,
	}

	if err = app.store.Targets.Create(ctx, &target); err != nil {
		switch err {
		case store.ErrDuplicate:
			app.conflictError(w, r, store.ErrDuplicateName)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, target); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAllTargets godoc
//
//	@Summary		Fetch all target
//	@Description	Fetch all target
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.Target
//	@Failure		403	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/targets [get]
func (app *application) fetchTargetsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targets, err := app.store.Targets.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, targets); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetTarget godoc
//
//	@Summary		Fetches a target
//	@Description	Fetches a target by ID
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Param			targetId	path		int	true	"Target ID"
//	@Success		200			{object}	store.Target
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/targets/{targetId} [get]
func (app *application) getTargetHandler(w http.ResponseWriter, r *http.Request) {
	target := getTargetFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, target); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeleteTarget godoc
//
//	@Summary		Deletes a target
//	@Description	Deletes a target by ID
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Target ID"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/targets/{targetId} [delete]
func (app *application) deleteTargetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	target := getTargetFromContext(r)

	if err := app.store.Targets.Delete(ctx, target.ID); err != nil {
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

type UpdateTargetPayload struct {
	Name       *string `json:"name" validate:"omitempty,max=40"`
	BodyPartID *int64  `json:"bodypart_id" validate:"omitempty"`
}

// UddateTarget godoc
//
//	@Summary		Update a target
//	@Description	Update a target by ID
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Param			targetId	path		int					true	"Target ID"
//	@Param			targetId	body		UpdateTargetPayload	true	"Target ID"
//	@Success		200			{object}	store.Target
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/targets/{targetId} [patch]
func (app *application) updateTargetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	target := getTargetFromContext(r)

	var payload UpdateTargetPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Name != nil {
		target.Name = *payload.Name
	}

	if payload.BodyPartID != nil {
		target.BodyPartID = *payload.BodyPartID
	}

	if err := app.store.Targets.Update(ctx, target); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, target); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) targetContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParam(r, "targetId")

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			app.badRequest(w, r, errors.New("invalid target id"))
			return
		}

		target, err := app.store.Targets.GetByID(ctx, intId)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFound(w, r)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, targetCtxKey, target)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTargetFromContext(r *http.Request) *store.PresentableTarget {
	target, _ := r.Context().Value(targetCtxKey).(*store.PresentableTarget)
	return target
}
