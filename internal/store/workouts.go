package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type WorkoutStore struct {
	db *sql.DB
}

type WorkoutDifficulty string

type CreateWorkoutPayload struct {
	Name             string            `json:"name"`
	BodyPartID       int64             `json:"bodypart_id"`
	EquipmentID      int64             `json:"equipment_id"`
	GifUrl           string            `json:"gif_url"`
	Instructions     []string          `json:"instructions"`
	CaloriesBurned   uint8             `json:"calories_burned"`
	DurationMinutes  uint8             `json:"duration_minutes"`
	Difficulty       WorkoutDifficulty `json:"difficulty"`
	PrimaryTarget    int64             `json:"primary_target"`
	SecondaryTargets []int64           `json:"secondary_targets"`
}
type Workout struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	BodyPartID      int64    `json:"bodypart_id"`
	EquipmentID     int64    `json:"equipment_id"`
	GifUrl          string   `json:"gif_url"`
	Instructions    []string `json:"instructions"`
	CaloriesBurned  uint8    `json:"calories_burned"`
	DurationMinutes uint8    `json:"duration_minutes"`
	Difficulty      string   `json:"difficulty"`
}

type PresentableWorkout struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	GifUrl           string    `json:"gif_url"`
	Instructions     []string  `json:"instructions"`
	CaloriesBurned   uint8     `json:"calories_burned"`
	DurationMinutes  uint8     `json:"duration_minutes"`
	Difficulty       string    `json:"difficulty"`
	BodyPart         string    `json:"body_part"`
	Equipment        string    `json:"equipment"`
	PrimaryTarget    string    `json:"primary_target"`
	SecondaryTargets []*string `json:"secondary_targets"`
}

func (s *WorkoutStore) create(ctx context.Context, tx *sql.Tx, workout *Workout) error {
	query := `
    INSERT INTO workout
    (name, bodypart_id, equipment_id, gif_url, instructions, calories_burned, duration_minutes, difficulty)
    VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id
    ;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := tx.QueryRowContext(
		ctx,
		query,
		&workout.Name,
		&workout.BodyPartID,
		&workout.EquipmentID,
		&workout.GifUrl,
		pq.Array(&workout.Instructions),
		&workout.CaloriesBurned,
		&workout.DurationMinutes,
		&workout.Difficulty,
	).Scan(&workout.ID); err != nil {
		// check unique constraints validation
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrDuplicate
		}
		return err
	}

	return nil
}

func (s *WorkoutStore) CreateAndLinkTargets(ctx context.Context, workout *Workout, primaryTargetId int64, secondaryTargetIds []int64) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		// create workout
		if err := s.create(ctx, tx, workout); err != nil {
			return err
		}

		// create links between workout and targets
		if err := s.linkTargets(ctx, tx, workout.ID, primaryTargetId, secondaryTargetIds); err != nil {
			return err
		}

		return nil
	})
}

func (s *WorkoutStore) linkTargets(ctx context.Context, tx *sql.Tx, workoutId int64, primaryTargetId int64, secondaryTargetIds []int64) error {
	linkPrimaryTargetquery := `
    INSERT INTO workout_target (workout_id, target_id, type)
    VALUES ($1, $2, 'primary')
    ;`

	linkSecondaryTargetquery := `
    INSERT INTO workout_target (workout_id, target_id, type)
    VALUES ($1, $2, 'secondary')
    ;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, linkPrimaryTargetquery, workoutId, primaryTargetId)
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

	for _, targetId := range secondaryTargetIds {
		res, err := tx.ExecContext(ctx, linkSecondaryTargetquery, workoutId, targetId)
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
	}

	return nil

}

func (s *WorkoutStore) GetAll(ctx context.Context) ([]PresentableWorkout, error) {
	query := `
    SELECT
    w.id, w.name, b.name, e.name, w.gif_url, w.difficulty, w.instructions, w.calories_burned, w.duration_minutes
    FROM workout w
    JOIN body_part b ON w.bodypart_id = b.id
    LEFT JOIN equipment e ON w.equipment_id = e.id
    ;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presentableWorkouts []PresentableWorkout
	for rows.Next() {
		var p PresentableWorkout
		err := rows.Scan(

			&p.ID,
			&p.Name,
			&p.BodyPart,
			&p.Equipment,
			&p.GifUrl,
			&p.Difficulty,
			pq.Array(&p.Instructions),
			&p.CaloriesBurned,
			&p.DurationMinutes,
		)
		if err != nil {
			return nil, err
		}

		primaryTarget, secondaryTargets, err := GetTargetsByWorkoutID(s.db, ctx, p.ID)
		p.PrimaryTarget = *primaryTarget
		p.SecondaryTargets = secondaryTargets
		if err != nil {
			return nil, err
		}
		presentableWorkouts = append(presentableWorkouts, p)
	}

	return presentableWorkouts, nil
}

func (s *WorkoutStore) GetByID(ctx context.Context, id int64) (*PresentableWorkout, error) {
	var presentableWorkout PresentableWorkout

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
    SELECT
    w.id, w.name, b.name, e.name, w.gif_url, w.difficulty, w.instructions, w.calories_burned, w.duration_minutes
    FROM workout w
    JOIN body_part b ON w.bodypart_id = b.id
    LEFT JOIN equipment e ON w.equipment_id = e.id
    WHERE w.id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&presentableWorkout.ID,
		&presentableWorkout.Name,
		&presentableWorkout.BodyPart,
		&presentableWorkout.Equipment,
		&presentableWorkout.GifUrl,
		&presentableWorkout.Difficulty,
		pq.Array(&presentableWorkout.Instructions),
		&presentableWorkout.CaloriesBurned,
		&presentableWorkout.DurationMinutes,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	primaryTarget, secondaryTargets, err := GetTargetsByWorkoutID(s.db, ctx, presentableWorkout.ID)
	presentableWorkout.PrimaryTarget = *primaryTarget
	presentableWorkout.SecondaryTargets = secondaryTargets
	if err != nil {
		return nil, err
	}

	return &presentableWorkout, nil
}

func (s *WorkoutStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM workout WHERE id = $1;`

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

func (s *WorkoutStore) Update(ctx context.Context, workout *PresentableWorkout) error {
	query := `UPDATE workout SET name = $1, bodypart_id = $2 WHERE id = $3;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, workout.Name, workout.ID)
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
