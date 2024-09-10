package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
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

	//go:embed sql/label_update.sql
	updateLabelQuery string

	//go:embed sql/label_delete.sql
	deleteLabelQuery string
)

func (db *DB) CreateLabel(label *model.Label) (*model.Label, error) {
	return label, db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		return label.ToSyncStatus(model.CreateOperation), tx.QueryRow(
			createLabelQuery,
			label.UserID,
			label.Name,
			label.Color,
		).Scan(&label.ID, &label.CreatedAt, &label.UpdatedAt)
	})
}

func (db *DB) GetLabelByID(userID, id int64) (*model.Label, error) {
	label := new(model.Label)

	err := db.QueryRow(getLabelByIDQuery, id, userID).Scan(
		&label.ID,
		&label.UserID,
		&label.Name,
		&label.Color,
		&label.CreatedAt,
		&label.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	return label, nil
}

func (db *DB) GetLabelByName(userID int64, name string) (*model.Label, error) {
	label := new(model.Label)

	err := db.QueryRow(getLabelByNameQuery, userID, name).Scan(
		&label.ID,
		&label.UserID,
		&label.Name,
		&label.Color,
		&label.CreatedAt,
		&label.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	return label, nil
}

func (db *DB) GetLabels(userID int64) ([]*model.Label, error) {
	var labels []*model.Label

	rows, err := db.Query(getLabelsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels: %w", err)
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
	return label, db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		err := tx.QueryRow(
			updateLabelQuery,
			label.ID,
			label.UserID,
			label.Name,
			label.Color,
		).Scan(&label.CreatedAt, &label.UpdatedAt)

		return label.ToSyncStatus(model.UpdateOperation), err
	})
}

func (db *DB) DeleteLabel(label *model.Label) error {
	return db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		if _, err := db.Exec(deleteLabelQuery, label.ID, label.UserID); err != nil {
			return nil, fmt.Errorf("failed to delete label: %w", err)
		}

		return label.ToSyncStatus(model.DeleteOperation), nil
	})
}
