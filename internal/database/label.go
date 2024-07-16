package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) LabelExistsByName(user int64, name string) bool {
	var result bool
	query := `SELECT true FROM labels WHERE user_id = $1 AND name = $2`

	_ = db.QueryRow(query, user, name).Scan(&result)
	return result
}

func (db *DB) CreateLabel(user int64, label *model.Label) (*model.Label, error) {
	if db.LabelExistsByName(user, label.Name) {
		return nil, fmt.Errorf("database: label already exists")
	}

	query := `
		INSERT INTO labels 
			(user_id, name, color)
		VALUES 
			($1, $2, $3)
		RETURNING 
			id, created_at, updated_at
	`
	err := db.QueryRow(query, user, label.Name, label.Color).Scan(&label.ID, &label.CreatedAt, &label.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create label: %v", err)
	}

	return label, nil
}

func (db *DB) GetLabelByID(userID, id int64) (*model.Label, error) {
	label := new(model.Label)
	query := `
		SELECT
			id,
			user_id,
			name,
			color,
			created_at,
			updated_at
		FROM
			labels
		WHERE 
			user_id = $1 AND id = $2
	`

	err := db.QueryRow(query, userID, id).Scan(&label.ID, &label.UserID, &label.Name, &label.Color, &label.CreatedAt, &label.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get label: %v", err)
	}

	return label, nil
}

func (db *DB) GetLabelByName(userID int64, name string) (*model.Label, error) {
	label := new(model.Label)
	query := `
		SELECT
			id,
			user_id,
			name,
			color,
			created_at,
			updated_at
		FROM
			labels
		WHERE 
			user_id = $1 AND name = $2
	`

	err := db.QueryRow(query, userID, name).Scan(&label.ID, &label.UserID, &label.Name, &label.Color, &label.CreatedAt, &label.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get label: %v", err)
	}

	return label, nil
}

func (db *DB) GetLabels(userID int64) ([]*model.Label, error) {
	var labels []*model.Label
	query := `
		SELECT
			id,
			user_id,
			name,
			color,
			created_at,
			updated_at
		FROM 
			labels
		WHERE 
			user_id = $1
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get labels: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		l := &model.Label{}

		if err := rows.Scan(&l.ID, &l.UserID, &l.Name, &l.Color, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("database: unable to get labels: %v", err)
		}

		labels = append(labels, l)
	}

	return labels, nil
}

func (db *DB) UpdateLabel(label *model.Label) (*model.Label, error) {
	query := `
		UPDATE 
			labels
		SET 
			name = $3,
			color = $4,
			updated_at = NOW()
		WHERE 
			id = $1 AND user_id = $2
		RETURNING 
			created_at,
			updated_at
	`

	err := db.QueryRow(
		query,
		label.ID,
		label.UserID,
		label.Name,
		label.Color,
	).Scan(&label.CreatedAt, &label.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to update label: %v", err)
	}

	return label, nil
}

func (db *DB) DeleteLabel(userID, id int64) error {
	query := `DELETE FROM labels WHERE id = $1 AND user_id = $2`

	_, err := db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("database: unable to delete label: %v", err)
	}

	return nil
}
