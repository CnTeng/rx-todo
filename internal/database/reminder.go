package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

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

	//go:embed sql/reminder_get_by_updated_at.sql
	getRemindersByUpdateTimeQuery string

	//go:embed sql/reminder_update.sql
	updateReminderQuery string

	//go:embed sql/reminder_delete.sql
	deleteReminderQuery string
)

func (db *DB) CreateReminder(reminder *model.Reminder) (*model.Reminder, error) {
	if err := db.QueryRow(
		createReminderQuery,
		reminder.UserID,
		reminder.TaskID,
		reminder.Due.Date,
		reminder.Due.Recurring,
	).Scan(
		&reminder.ID,
		&reminder.CreatedAt,
		&reminder.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return reminder, nil
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

func (db *DB) GetRemindersByUpdateTime(userID int64, updateTime *time.Time) ([]*model.Reminder, error) {
	var reminders []*model.Reminder

	rows, err := db.Query(getRemindersByUpdateTimeQuery, userID, updateTime)
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
	if err := db.QueryRow(
		updateReminderQuery,
		reminder.ID,
		reminder.UserID,
		reminder.Due.Date,
		reminder.Due.Recurring,
	).Scan(&reminder.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to update reminder: %w", err)
	}

	return reminder, nil
}

func (db *DB) DeleteReminder(reminder *model.Reminder) error {
	return db.withTx(func(tx *sql.Tx) error {
		if _, err := tx.Exec(deleteReminderQuery, reminder.ID, reminder.UserID); err != nil {
			return fmt.Errorf("failed to delete reminder: %w", err)
		}

		if _, err := tx.Exec(createDeletionLogQuery, reminder.UserID, "reminder", reminder.ID); err != nil {
			return fmt.Errorf("failed to create deletion log: %w", err)
		}
		return nil
	})
}
