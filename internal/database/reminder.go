package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

var (
	//go:embed sql/reminder_create.sql
	createReminderQuery string

	//go:embed sql/reminder_get_by_id.sql
	getReminderByIDQuery string

	//go:embed sql/reminder_get_by_task.sql
	getReminderByTaskIDQuery string

	//go:embed sql/reminder_get_all.sql
	getRemindersQuery string

	//go:embed sql/reminder_update.sql
	updateReminderQuery string

	//go:embed sql/reminder_delete.sql
	deleteReminderQuery string
)

func (db *DB) CreateReminder(reminder *model.Reminder) (*model.Reminder, error) {
	return reminder, db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		err := tx.QueryRow(
			createReminderQuery,
			reminder.UserID,
			reminder.TaskID,
			reminder.Due.Date,
			reminder.Due.Recurring,
		).Scan(&reminder.ID, &reminder.CreatedAt, &reminder.UpdatedAt)

		return reminder.ToSyncStatus(model.CreateOperation), err
	})
}

func (db *DB) GetReminderByID(userID, id int64) (*model.Reminder, error) {
	reminder := new(model.Reminder)

	err := db.QueryRow(getReminderByIDQuery, id, userID).Scan(
		&reminder.ID,
		&reminder.UserID,
		&reminder.TaskID,
		&reminder.Due.Date,
		&reminder.Due.Recurring,
		&reminder.CreatedAt,
		&reminder.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder: %w", err)
	}

	return reminder, nil
}

func (db *DB) GetReminderByTaskID(userID, taskID int64) (*model.Reminder, error) {
	reminder := new(model.Reminder)

	err := db.QueryRow(getReminderByTaskIDQuery, userID, taskID).Scan(
		&reminder.ID,
		&reminder.UserID,
		&reminder.TaskID,
		&reminder.Due.Date,
		&reminder.Due.Recurring,
		&reminder.CreatedAt,
		&reminder.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder: %w", err)
	}

	return reminder, nil
}

func (db *DB) GetReminders(userID int64) ([]*model.Reminder, error) {
	var reminders []*model.Reminder

	rows, err := db.Query(getRemindersQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		reminder := &model.Reminder{}
		if err := rows.Scan(
			&reminder.ID,
			&reminder.UserID,
			&reminder.TaskID,
			&reminder.Due.Date,
			&reminder.Due.Recurring,
			&reminder.CreatedAt,
			&reminder.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("database: unable to get reminder: %v", err)
		}

		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

func (db *DB) UpdateReminder(reminder *model.Reminder) (*model.Reminder, error) {
	return reminder, db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		err := db.QueryRow(
			updateReminderQuery,
			reminder.ID,
			reminder.UserID,
			reminder.Due.Date,
			reminder.Due.Recurring,
		).Scan(&reminder.UpdatedAt)

		return reminder.ToSyncStatus(model.UpdateOperation), err
	})
}

func (db *DB) DeleteReminder(reminder *model.Reminder) error {
	return db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		_, err := db.Exec(deleteReminderQuery, reminder.ID, reminder.UserID)
		return reminder.ToSyncStatus(model.DeleteOperation), err
	})
}
