package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type TargetStore struct {
	db *sql.DB
}

type Target struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	BodyPartID int64  `json:"bodypart_id"`
}

type PresentableTarget struct {
	Target
	BodyPart string `json:"body_part"`
}

func (s *TargetStore) Create(ctx context.Context, target *Target) error {
	query := `INSERT INTO target (name, bodypart_id) VALUES ($1, $2) RETURNING id;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, &target.Name, &target.BodyPartID).Scan(&target.ID); err != nil {
		// check unique constraints validation
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrDuplicate
		}
		return err
	}

	return nil
}

func (s *TargetStore) GetAll(ctx context.Context) ([]PresentableTarget, error) {
	query := `
    SELECT
    t.id, t.name, b.id, b.name
    FROM target t
    JOIN body_part b on t.bodypart_id = b.id
    ;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presentableTargets []PresentableTarget
	for rows.Next() {
		var t PresentableTarget
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.BodyPartID,
			&t.BodyPart,
		)
		if err != nil {
			return nil, err
		}
		presentableTargets = append(presentableTargets, t)
	}

	return presentableTargets, nil
}

func (s *TargetStore) GetByID(ctx context.Context, id int64) (*PresentableTarget, error) {
	var presentableTarget PresentableTarget

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
    SELECT
    t.id, t.name, b.id, b.name
    FROM target t
    JOIN body_part b on t.bodypart_id = b.id
    WHERE t.id = $1`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&presentableTarget.ID, &presentableTarget.Name, &presentableTarget.BodyPartID, &presentableTarget.BodyPart)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &presentableTarget, nil
}

func (s *TargetStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM target WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *TargetStore) Update(ctx context.Context, target *PresentableTarget) error {
	query := `UPDATE target SET name = $1, bodypart_id = $2 WHERE id = $3;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, target.Name, target.BodyPartID, target.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil

}

func GetTargetsByWorkoutID(db *sql.DB, ctx context.Context, workoutId int64) (*string, []*string, error) {
	var primaryTarget *string
	var secondaryTargets []*string

	query := `
    SELECT t.name, wt.type
    FROM workout_target wt
    JOIN target t ON t.id = wt.target_id
    WHERE wt.workout_id = $1`

	rows, err := db.QueryContext(ctx, query, workoutId)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Target
		var targetType *string
		if err := rows.Scan(&t.Name, &targetType); err != nil {
			return nil, nil, err
		}
		if targetType != nil && *targetType == "primary" {
			primaryTarget = &t.Name
		} else {
			secondaryTargets = append(secondaryTargets, &t.Name)
		}
	}

	return primaryTarget, secondaryTargets, nil
}
