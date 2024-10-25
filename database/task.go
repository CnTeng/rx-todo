package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/CnTeng/rx-todo/model"
)

var (
	//go:embed sql/task_create.sql
	createTaskQuery string

	//go:embed sql/task_labels_create.sql
	createTaskLabelsQuery string

	//go:embed sql/task_labels_delete.sql
	deleteTaskLabelsQuery string

	//go:embed sql/task_get_by_id.sql
	getTaskByIDQuery string

	//go:embed sql/task_get_all.sql
	getTasksQuery string

	//go:embed sql/task_get_by_updated_at.sql
	getTasksByUpdateTimeQuery string

	//go:embed sql/task_get_new_position.sql
	getTaskNewPositionQuery string

	//go:embed sql/task_get_top_position.sql
	getTaskTopPositionQuery string

	//go:embed sql/task_update.sql
	updateTaskQuery string

	//go:embed sql/task_update_position.sql
	updateTaskPositionQuery string

	//go:embed sql/task_update_done.sql
	updateTaskDoneQuery string

	//go:embed sql/task_update_archived.sql
	updateTaskArchivedQuery string

	//go:embed sql/task_delete.sql
	deleteTaskQuery string
)

func (db *DB) CreateTask(task *model.Task) (*model.Task, error) {
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

	return task, db.withTx(func(tx *sql.Tx) error {
		if err := tx.QueryRow(
			createTaskQuery,
			task.UserID,
			task.Name,
			task.Description,
			dueDate,
			dueRecurring,
			durationAmount,
			durationUnit,
			task.Priority,
			task.ProjectID,
			task.ParentID,
			task.Position,
		).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		for _, label := range task.Labels {
			if _, err := tx.Exec(
				createTaskLabelsQuery,
				task.UserID,
				label.Name,
				task.ID,
			); err != nil {
				return fmt.Errorf("failed to create task labels: %w", err)
			}
		}

		return nil
	})
}

func (db *DB) GetTaskByID(userID, id int64) (*model.Task, error) {
	var dueData *time.Time
	var dueRecurring *bool
	var durationAmount *int
	var durationUnit *string

	task := &model.Task{}

	if err := db.QueryRow(
		getTaskByIDQuery,
		userID,
		id,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Name,
		&task.Description,
		&dueData,
		&dueRecurring,
		&durationAmount,
		&durationUnit,
		&task.Priority,
		&task.ProjectID,
		&task.ParentID,
		&task.Position,
		&task.Done,
		&task.DoneAt,
		&task.Archived,
		&task.ArchivedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if (dueData != nil) && (dueRecurring != nil) {
		task.Due = &model.Due{Date: dueData, Recurring: dueRecurring}
	}

	if (durationAmount != nil) && (durationUnit != nil) {
		task.Duration = &model.Duration{Amount: durationAmount, Unit: durationUnit}
	}

	labels, err := db.GetLabelsByTaskID(task.ID, task.UserID)
	if err != nil {
		return nil, err
	}
	task.Labels = labels

	return task, nil
}

func (db *DB) getTasks(query string, args ...any) ([]*model.Task, error) {
	tasks := []*model.Task{}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &model.Task{}

		var dueData *time.Time
		var dueRecurring *bool
		var durationAmount *int
		var durationUnit *string

		if err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Name,
			&task.Description,
			&dueData,
			&dueRecurring,
			&durationAmount,
			&durationUnit,
			&task.Priority,
			&task.ProjectID,
			&task.ParentID,
			&task.Position,
			&task.Done,
			&task.DoneAt,
			&task.Archived,
			&task.ArchivedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get tasks: %w", err)
		}

		if (dueData != nil) && (dueRecurring != nil) {
			task.Due = &model.Due{Date: dueData, Recurring: dueRecurring}
		}

		if (durationAmount != nil) && (durationUnit != nil) {
			task.Duration = &model.Duration{Amount: durationAmount, Unit: durationUnit}
		}

		labels, err := db.GetLabelsByTaskID(task.ID, task.UserID)
		if err != nil {
			return nil, err
		}
		task.Labels = labels

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *DB) GetTasks(user int64) ([]*model.Task, error) {
	return db.getTasks(getTasksQuery, user)
}

func (db *DB) GetTasksByUpdateTime(user int64, updateTime *time.Time) ([]*model.Task, error) {
	return db.getTasks(getTasksByUpdateTimeQuery, user, updateTime)
}

func (db *DB) UpdateTask(task *model.Task) (*model.Task, error) {
	return task, db.withTx(func(tx *sql.Tx) error {
		if err := tx.QueryRow(
			updateTaskQuery,
			task.ID,
			task.UserID,
			task.Name,
			task.Description,
			task.Due.Date,
			task.Due.Recurring,
			task.Duration.Amount,
			task.Duration.Unit,
			task.Priority,
		).Scan(&task.UpdatedAt); err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		if _, err := tx.Exec(deleteTaskLabelsQuery, task.ID); err != nil {
			return fmt.Errorf("failed to delete task labels: %w", err)
		}

		for _, label := range task.Labels {
			if _, err := tx.Exec(
				createTaskLabelsQuery,
				task.UserID,
				label.Name,
				task.ID,
			); err != nil {
				return fmt.Errorf("failed to create task labels: %w", err)
			}
		}

		return nil
	})
}

func (db *DB) UpdateTaskPosition(task *model.Task, previousID int64) (*model.Task, error) {
	return task, db.withTx(func(tx *sql.Tx) error {
		if previousID == 0 {
			if err := tx.QueryRow(
				getTaskTopPositionQuery,
				task.UserID,
			).Scan(&task.Position); err != nil {
				return fmt.Errorf("failed to get project top position: %w", err)
			}
		} else {
			if err := tx.QueryRow(
				getTaskNewPositionQuery,
				previousID,
				task.UserID,
				task.ProjectID,
				task.ParentID,
			).Scan(&task.Position); err != nil {
				return fmt.Errorf("failed to get project new position: %w", err)
			}
		}

		if err := tx.QueryRow(
			updateTaskPositionQuery,
			task.ID,
			task.UserID,
			task.ProjectID,
			task.ParentID,
			task.Position,
		).Scan(&task.UpdatedAt); err != nil {
			return fmt.Errorf("failed to update project position: %w", err)
		}

		return nil
	})
}

func (db *DB) UpdateTaskDoneStatus(task *model.Task) (*model.Task, error) {
	if err := db.QueryRow(
		updateTaskDoneQuery,
		task.ID,
		task.UserID,
		task.Done,
	).Scan(&task.DoneAt); err != nil {
		return nil, fmt.Errorf("failed to update project done status: %w", err)
	}
	return task, nil
}

func (db *DB) UpdateTaskArchivedStatus(task *model.Task) (*model.Task, error) {
	if err := db.QueryRow(
		updateTaskArchivedQuery,
		task.ID,
		task.UserID,
		task.Archived,
	).Scan(&task.ArchivedAt); err != nil {
		return nil, fmt.Errorf("failed to update project archived status: %w", err)
	}
	return task, nil
}

func (db *DB) DeleteTask(task *model.Task) error {
	return db.withTx(func(tx *sql.Tx) error {
		if _, err := db.Exec(deleteTaskQuery, task.ID, task.UserID); err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		if _, err := tx.Exec(createDeletionLogQuery, task.UserID, "task", task.ID); err != nil {
			return fmt.Errorf("failed to create deletion log: %w", err)
		}
		return nil
	})
}
