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
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		BodyParts: &BodyPartStore{db},
		Targets:   &TargetStore{db},
	}
}
