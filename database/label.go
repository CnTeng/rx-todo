package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/CnTeng/rx-todo/model"
)

var (
	//go:embed sql/label_create.sql
	createLabelQuery string

	//go:embed sql/label_get_by_id.sql
	getLabelByIDQuery string

	//go:embed sql/label_get_by_name.sql
	getLabelByNameQuery string

	//go:embed sql/label_get_all.sql
	getLabelsQuery string

	//go:embed sql/label_get_by_updated_at.sql
	getLabelsByUpdateTimeQuery string

	//go:embed sql/label_get_by_task_id.sql
	getLabelsByTaskIDQuery string

	//go:embed sql/label_update.sql
	updateLabelQuery string

	//go:embed sql/label_delete.sql
	deleteLabelQuery string
)

func (db *DB) CreateLabel(label *model.Label) (*model.Label, error) {
	if err := db.QueryRow(
		createLabelQuery,
		label.UserID,
		label.Name,
		label.Color,
	).Scan(&label.ID, &label.CreatedAt, &label.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	return label, nil
}

func (db *DB) getLabel(query string, args ...any) (*model.Label, error) {
	label := &model.Label{}

	if err := db.QueryRow(query, args).Scan(
		&label.ID,
		&label.UserID,
		&label.Name,
		&label.Color,
		&label.CreatedAt,
		&label.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	return label, nil
}

func (db *DB) GetLabelByID(id, userID int64) (*model.Label, error) {
	return db.getLabel(getLabelByIDQuery, id, userID)
}

func (db *DB) GetLabelByName(name string, userID int64) (*model.Label, error) {
	return db.getLabel(getLabelByNameQuery, userID, name)
}

func (db *DB) getLabels(query string, args ...any) ([]*model.Label, error) {
	labels := []*model.Label{}

	rows, err := db.Query(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		l := &model.Label{}

		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.Name,
			&l.Color,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("database: unable to get labels: %v", err)
		}

		labels = append(labels, l)
	}

	return labels, nil
}

func (db *DB) GetLabels(userID int64) ([]*model.Label, error) {
	return db.getLabels(getLabelsQuery, userID)
}

func (db *DB) GetLabelsByUpdateTime(updateTime *time.Time, userID int64) ([]*model.Label, error) {
	if updateTime == nil {
		return db.GetLabels(userID)
	}
	return db.getLabels(getLabelsByUpdateTimeQuery, userID, updateTime)
}

func (db *DB) GetLabelsByTaskID(taskID, userID int64) ([]*model.Label, error) {
	labels := []*model.Label{}

	rows, err := db.Query(getLabelsByTaskIDQuery, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task labels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		label := &model.Label{}
		if err := rows.Scan(&label.ID); err != nil {
			return nil, fmt.Errorf("failed to get task labels: %w", err)
		}

		if err := db.QueryRow(getLabelByIDQuery, label.ID, userID).Scan(
			&label.ID,
			&label.UserID,
			&label.Name,
			&label.Color,
			&label.CreatedAt,
			&label.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get label: %w", err)
		}

		labels = append(labels, label)
	}

	return labels, nil
}

func (db *DB) UpdateLabel(label *model.Label) (*model.Label, error) {
	if err := db.QueryRow(
		updateLabelQuery,
		label.ID,
		label.UserID,
		label.Name,
		label.Color,
	).Scan(&label.CreatedAt, &label.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to update label: %w", err)
	}

	return label, nil
}

func (db *DB) DeleteLabel(id, userID int64) error {
	return db.withTx(func(tx *sql.Tx) error {
		if _, err := tx.Exec(deleteLabelQuery, id, userID); err != nil {
			return fmt.Errorf("failed to delete label: %w", err)
		}

		if _, err := tx.Exec(createDeletionLogQuery, userID, "label", id); err != nil {
			return fmt.Errorf("failed to create deletion log: %w", err)
		}
		return nil
	})
}
