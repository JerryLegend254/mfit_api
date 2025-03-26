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

type workoutContextKey string

var workoutCtxKey workoutContextKey = "workout"

type CreateWorkoutPayload struct {
	Name             string   `json:"name" validate:"required,max=40"`
	BodyPartID       int64    `json:"bodypart_id" validate:"required"`
	EquipmentID      int64    `json:"equipment_id" validate:"required"`
	GifUrl           string   `json:"gif_url"`
	Instructions     []string `json:"instructions"`
	CaloriesBurned   uint8    `json:"calories_burned"`
	DurationMinutes  uint8    `json:"duration_minutes"`
	Difficulty       string   `json:"difficulty" validate:"required"`
	PrimaryTarget    int64    `json:"primary_target" validate:"required"`
	SecondaryTargets []int64  `json:"secondary_targets" validate:"required"`
}

// CreateWorkout godoc
//
//	@Summary		Creates a workout
//	@Description	Creates a workout
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateWorkoutPayload	true	"Workout payload"
//	@Success		201		{object}	store.Workout
//	@Failure		400		{object}	error
//	@Failure		403		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workouts [post]
func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	workout := store.Workout{
		Name:            payload.Name,
		BodyPartID:      payload.BodyPartID,
		EquipmentID:     payload.EquipmentID,
		GifUrl:          payload.GifUrl,
		Instructions:    payload.Instructions,
		CaloriesBurned:  payload.CaloriesBurned,
		DurationMinutes: payload.DurationMinutes,
		Difficulty:      payload.Difficulty,
	}

	if err = app.store.Workouts.CreateAndLinkTargets(ctx, &workout, payload.PrimaryTarget, payload.SecondaryTargets); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusCreated, &workout); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAllWorkouts godoc
//
//	@Summary		Fetch all workout
//	@Description	Fetch all workout
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.Workout
//	@Failure		403	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workouts [get]

func (app *application) fetchWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workouts, err := app.store.Workouts.GetAll(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, workouts); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetWorkout godoc
//
//	@Summary		Fetches a workout
//	@Description	Fetches a workout by ID
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			workoutId	path		int	true	"Workout ID"
//	@Success		200			{object}	store.Workout
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workouts/{workoutId} [get]
func (app *application) getWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	workout := getWorkoutFromContext(r)
	if err := app.jsonResponse(w, http.StatusOK, workout); err != nil {
		app.internalServerError(w, r, err)
	}
}

// // DeleteWorkout godoc
// //
// //	@Summary		Deletes a workout
// //	@Description	Deletes a workout by ID
// //	@Tags			workouts
// //	@Accept			json
// //	@Produce		json
// //	@Param			id	path	int	true	"Workout ID"
// //	@Success		204
// //	@Failure		400	{object}	error
// //	@Failure		401	{object}	error
// //	@Failure		404	{object}	error
// //	@Failure		500	{object}	error
// //	@Security		ApiKeyAuth
// //	@Router			/workouts/{workoutId} [delete]
//
//	func (app *application) deleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		workout := getWorkoutFromContext(r)
//
//		if err := app.store.Workouts.Delete(ctx, workout.ID); err != nil {
//			switch err {
//			case store.ErrNotFound:
//				app.notFound(w, r)
//			default:
//				app.internalServerError(w, r, err)
//			}
//			return
//		}
//
//		w.WriteHeader(http.StatusNoContent)
//	}
//
//	type UpdateWorkoutPayload struct {
//		Name       *string `json:"name" validate:"omitempty,max=40"`
//		BodyPartID *int64  `json:"bodypart_id" validate:"omitempty"`
//	}
//
// // UddateWorkout godoc
// //
// //	@Summary		Update a workout
// //	@Description	Update a workout by ID
// //	@Tags			workouts
// //	@Accept			json
// //	@Produce		json
// //	@Param			workoutId	path		int						true	"Workout ID"
// //	@Param			workoutId	body		UpdateWorkoutPayload	true	"Workout ID"
// //	@Success		200			{object}	store.Workout
// //	@Failure		400			{object}	error
// //	@Failure		401			{object}	error
// //	@Failure		404			{object}	error
// //	@Failure		500			{object}	error
// //	@Security		ApiKeyAuth
// //	@Router			/workouts/{workoutId} [patch]
//
//	func (app *application) updateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//
//		workout := getWorkoutFromContext(r)
//
//		var payload UpdateWorkoutPayload
//
//		if err := readJSON(w, r, &payload); err != nil {
//			app.badRequest(w, r, err)
//			return
//		}
//
//		if err := Validate.Struct(payload); err != nil {
//			app.badRequest(w, r, err)
//			return
//		}
//
//		if payload.Name != nil {
//			workout.Name = *payload.Name
//		}
//
//		if payload.BodyPartID != nil {
//			workout.BodyPartID = *payload.BodyPartID
//		}
//
//		if err := app.store.Workouts.Update(ctx, workout); err != nil {
//			app.internalServerError(w, r, err)
//			return
//		}
//
//		if err := app.jsonResponse(w, http.StatusOK, workout); err != nil {
//			app.internalServerError(w, r, err)
//			return
//		}
//
// }
func (app *application) workoutContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParam(r, "workoutId")

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			app.badRequest(w, r, errors.New("invalid workout id"))
			return
		}

		workout, err := app.store.Workouts.GetByID(ctx, intId)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFound(w, r)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, workoutCtxKey, workout)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getWorkoutFromContext(r *http.Request) *store.PresentableWorkout {
	workout, _ := r.Context().Value(workoutCtxKey).(*store.PresentableWorkout)
	return workout
}
