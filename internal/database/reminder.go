package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) CreateReminder(reminder *model.Reminder) (*model.Reminder, error) {
	query := `
    INSERT INTO reminders 
      (user_id, task_id, due)
    VALUES 
      ($1, $2, ROW($3, $4))
    RETURNING 
      id, created_at, updated_at
  `

	err := db.QueryRow(
		query,
		reminder.UserID,
		reminder.TaskID,
		reminder.Due.Date,
		reminder.Due.Recurring,
	).Scan(&reminder.ID, &reminder.CreatedAt, &reminder.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create reminder: %v", err)
	}

	return reminder, nil
}

func (db *DB) GetReminderByID(userID, id int64) (*model.Reminder, error) {
	reminder := new(model.Reminder)
	query := `
    SELECT
      id,
      user_id,
      task_id,
      (due).date,
      (due).recurring,
      created_at,
      updated_at
    FROM
      reminders
    WHERE 
      id = $1 AND user_id = $2
  `

	err := db.QueryRow(query, id, userID).Scan(
		&reminder.ID,
		&reminder.UserID,
		&reminder.TaskID,
		&reminder.Due.Date,
		&reminder.Due.Recurring,
		&reminder.CreatedAt,
		&reminder.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get reminder: %v", err)
	}

	return reminder, nil
}

func (db *DB) GetReminderByTaskID(userID, taskID int64) (*model.Reminder, error) {
	reminder := new(model.Reminder)
	query := `
    SELECT
      id,
      user_id,
      task_id,
      (due).date,
      (due).recurring,
      created_at,
      updated_at
    FROM
      reminders
    WHERE 
      user_id = $1 AND task_id = $2
  `

	err := db.QueryRow(query, userID, taskID).Scan(
		&reminder.ID,
		&reminder.UserID,
		&reminder.TaskID,
		&reminder.Due.Date,
		&reminder.Due.Recurring,
		&reminder.CreatedAt,
		&reminder.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get reminder: %v", err)
	}

	return reminder, nil
}

func (db *DB) GetReminders(userID int64) ([]*model.Reminder, error) {
	var reminders []*model.Reminder

	query := `
    SELECT
      id,
      user_id,
      task_id,
      (due).date,
      (due).recurring,
      created_at,
      updated_at
    FROM
      reminders
    WHERE 
      user_id = $1
  `

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get tasks: %v", err)
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
	query := `
    UPDATE 
      reminders
    SET
      due = ROW($3, $4),
      updated_at = NOW()
    WHERE
      id = $1 AND user_id = $2
    RETURNING
      updated_at
  `

	err := db.QueryRow(
		query,
		reminder.ID,
		reminder.UserID,
		reminder.Due.Date,
		reminder.Due.Recurring,
	).Scan(&reminder.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to update reminder: %v", err)
	}

	return reminder, nil
}

func (db *DB) DeleteReminder(userID, id int64) error {
	query := `DELETE FROM reminders WHERE id = $1 AND user_id = $2`

	_, err := db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("database: unable to delete reminder: %v", err)
	}

	return nil
}
