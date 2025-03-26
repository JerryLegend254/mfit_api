package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
	ErrNotFound          = errors.New("resource not found")
)

type Storage struct {
	BodyParts interface {
		Create(context.Context, *BodyPart) error
		GetByID(context.Context, int64) (*BodyPart, error)
		GetAll(context.Context) ([]BodyPart, error)
		Update(context.Context, *BodyPart) error
		Delete(context.Context, int64) error
	}
	Targets interface {
		Create(context.Context, *Target) error
		GetByID(context.Context, int64) (*PresentableTarget, error)
		GetAll(context.Context) ([]PresentableTarget, error)
		Update(context.Context, *PresentableTarget) error
		Delete(context.Context, int64) error
	}
	Equipment interface {
		Create(context.Context, *Equipment) error
		GetByID(context.Context, int64) (*Equipment, error)
		GetAll(context.Context) ([]Equipment, error)
		Update(context.Context, *Equipment) error
		Delete(context.Context, int64) error
	}
	Workouts interface {
		CreateAndLinkTargets(context.Context, *Workout, int64, []int64) error
		GetByID(context.Context, int64) (*PresentableWorkout, error)
		GetAll(context.Context) ([]PresentableWorkout, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		BodyParts: &BodyPartStore{db},
		Targets:   &TargetStore{db},
		Equipment: &EquipmentStore{db},
		Workouts:  &WorkoutStore{db},
	}
}

func withTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
