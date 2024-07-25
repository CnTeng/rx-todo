package database

import (
	"fmt"
	"time"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) CreateTask(user int64, task *model.Task) (*model.Task, error) {
	var dueDate *time.Time
	var dueRecurring *bool
	var durationAmount *int
	var durationUnit *string

	if task.Due != nil {
		dueDate = task.Due.Date
		dueRecurring = task.Due.Recurring
	}

	if task.Duration != nil {
		durationAmount = task.Duration.Amount
		durationUnit = task.Duration.Unit
	}

	query := `
		INSERT INTO tasks
			(user_id, content, description, due, duration, priority, project_id, child_order)
		VALUES
			($1, $2, $3, ROW($4, $5), ROW($6, $7), $8, $9, $10)
		RETURNING
			id, created_at, updated_at
	`

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("database: unable to create task: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRow(
		query,
		user,
		task.Content,
		task.Description,
		dueDate,
		dueRecurring,
		durationAmount,
		durationUnit,
		task.Priority,
		task.ProjectID,
		task.ChildOrder,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create task: %v", err)
	}

	for _, labelName := range task.Labels {
		label, err := db.GetLabelByName(user, labelName)
		if err != nil {
			return nil, fmt.Errorf("database: label not found: %v", label)
		}

		_, err = tx.Exec(
			`
				INSERT INTO task_labels
					(task_id, label_id)
				VALUES
					($1, $2)
			`, task.ID, label.ID)
		if err != nil {
			return nil, fmt.Errorf("database: unable to create task labels: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("database: unable to commit transaction: %v", err)
	}

	return task, nil
}

func (db *DB) GetTaskByID(userID, id int64) (*model.Task, error) {
	var dueData *time.Time
	var dueRecurring *bool
	var durationAmount *int
	var durationUnit *string

	task := new(model.Task)
	query := `
		SELECT
			id,
			user_id,
			content,
			description,
			(due).date,
			(due).recurring,
			(duration).amount,
			(duration).unit,
			priority,
			project_id,
			parent_id,
			child_order,
			labels,
			done,
			done_at,
			archived,
			archived_at,
			created_at,
			updated_at
		FROM
			task_with_labels
		WHERE
			user_id = $1 AND id = $2
	`

	err := db.QueryRow(query, userID, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Content,
		&task.Description,
		&dueData,
		&dueRecurring,
		&durationAmount,
		&durationUnit,
		&task.Priority,
		&task.ProjectID,
		&task.ParentID,
		&task.ChildOrder,
		&task.Labels,
		&task.Done,
		&task.DoneAt,
		&task.Archived,
		&task.ArchivedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get task: %v", err)
	}

	if (dueData != nil) && (dueRecurring != nil) {
		task.Due = &model.Due{Date: dueData, Recurring: dueRecurring}
	}

	if (durationAmount != nil) && (durationUnit != nil) {
		task.Duration = &model.Duration{Amount: durationAmount, Unit: durationUnit}
	}

	return task, nil
}

func (db *DB) GetTasks(user int64) ([]*model.Task, error) {
	var tasks []*model.Task

	var dueData *time.Time
	var dueRecurring *bool
	var durationAmount *int
	var durationUnit *string

	query := `
		SELECT
			id,
			user_id,
			content,
			description,
			(due).date,
			(due).recurring,
			(duration).amount,
			(duration).unit,
			priority,
			project_id,
			parent_id,
			child_order,
			labels,
			done,
			done_at,
			archived,
			archived_at,
			created_at,
			updated_at
		FROM
			task_with_labels
		WHERE
			user_id = $1
	`

	rows, err := db.Query(query, user)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task

		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Content,
			&task.Description,
			&dueData,
			&dueRecurring,
			&durationAmount,
			&durationUnit,
			&task.Priority,
			&task.ProjectID,
			&task.ParentID,
			&task.ChildOrder,
			&task.Labels,
			&task.Done,
			&task.DoneAt,
			&task.Archived,
			&task.ArchivedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("database: unable to get tasks: %v", err)
		}

		if (dueData != nil) && (dueRecurring != nil) {
			task.Due = &model.Due{Date: dueData, Recurring: dueRecurring}
		}

		if (durationAmount != nil) && (durationUnit != nil) {
			task.Duration = &model.Duration{Amount: durationAmount, Unit: durationUnit}
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (db *DB) UpdateTask(task *model.Task) (*model.Task, error) {
	query := `
    UPDATE
			tasks
    SET
      content = $3,
      description = $4,
      due = $5,
      duration = $6,
      priority = $7,
			project_id = $8,
			parent_id = $9,
			child_order = $10,
			updated_at = now()
    WHERE
			id = $1 AND user_id = $2
		RETURNING
			updated_at
  `

	err := db.QueryRow(query, task.ID, task.UserID, task.Content, task.Description, task.Due, task.Duration, task.Priority, task.ProjectID, task.ParentID, task.ChildOrder).Scan(&task.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to update task: %v", err)
	}

	return task, nil
}

func (db *DB) DeleteTask(user, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	_, err := db.Exec(query, id, user)
	if err != nil {
		return fmt.Errorf("database: unable to delete tasks: %v", err)
	}

	return nil
}
