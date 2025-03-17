package store

import (
	"context"
	"database/sql"
)

type EquipmentStore struct {
	db *sql.DB
}

type Equipment struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *EquipmentStore) Create(ctx context.Context, equipment *Equipment) error {
	query := `INSERT INTO equipment (name) VALUES ($1)  RETURNING id;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRowContext(ctx, query, &equipment.Name).Scan(&equipment.ID); err != nil {
		return err
	}

	return nil
}

func (s *EquipmentStore) GetAll(ctx context.Context) ([]Equipment, error) {
	query := `
    SELECT
    id, name
    FROM equipment
    ;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipment []Equipment
	for rows.Next() {
		var e Equipment
		err := rows.Scan(
			&e.ID,
			&e.Name,
		)
		if err != nil {
			return nil, err
		}
		equipment = append(equipment, e)
	}

	return equipment, nil
}

func (s *EquipmentStore) GetByID(ctx context.Context, id int64) (*Equipment, error) {
	var equipment Equipment

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
    SELECT
    id, name
    FROM equipment
    WHERE id = $1`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&equipment.ID, &equipment.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &equipment, nil
}

func (s *EquipmentStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM equipment WHERE id = $1;`

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

func (s *EquipmentStore) Update(ctx context.Context, equipment *Equipment) error {
	query := `UPDATE equipment SET name = $1 WHERE id = $2;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, equipment.Name, equipment.ID)
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
