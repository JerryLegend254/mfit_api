package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type BodyPartStore struct {
	db *sql.DB
}

type BodyPart struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func (s *BodyPartStore) Create(ctx context.Context, bodyPart *BodyPart) error {
	query := `INSERT INTO body_part (name, image_url) VALUES ($1, $2) RETURNING id;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, &bodyPart.Name, &bodyPart.ImageUrl).Scan(&bodyPart.ID); err != nil {
		// check unique constraints validation
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrDuplicate
		}
		return err

	}

	return nil
}

func (s *BodyPartStore) GetAll(ctx context.Context) ([]BodyPart, error) {
	query := `SELECT id, name, image_url FROM body_part;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bodyParts []BodyPart
	for rows.Next() {
		var b BodyPart
		err := rows.Scan(
			&b.ID,
			&b.Name,
			&b.ImageUrl,
		)
		if err != nil {
			return nil, err
		}
		bodyParts = append(bodyParts, b)
	}

	return bodyParts, nil
}

func (s *BodyPartStore) GetByID(ctx context.Context, id int64) (*BodyPart, error) {
	var bodyPart BodyPart

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, name, image_url FROM body_part WHERE id = $1`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&bodyPart.ID, &bodyPart.Name, &bodyPart.ImageUrl)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &bodyPart, nil
}

func (s *BodyPartStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM body_part WHERE id = $1;`

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

func (s *BodyPartStore) Update(ctx context.Context, bodyPart *BodyPart) error {
	query := `UPDATE body_part SET name = $1, image_url = $2 WHERE id = $3;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, bodyPart.Name, bodyPart.ImageUrl, bodyPart.ID)
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
