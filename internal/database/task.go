package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) AddTask(user int64, r *model.TaskAddRequest) error {
	query := `
    INSERT INTO tasks (
      user_id,
      content,
      description,
      due,
      duration,
      priority,
      project_id,
      parent_id,
      child_order 
    )
    VALUES ($1, $2, $3, ROW($4, $5), $6, $7, $8, $9, $10)
  `

	_, err := db.Exec(query, user, r.Content, r.Description, r.Due.Data, r.Due.Recurring, r.Duration, r.Priority, r.ProjectID, r.ParentID, r.ChildOrder)
	if err != nil {
		return fmt.Errorf("database: unable to add task: %v", err)
	}

	return nil
}

func (db *DB) UpdateTask(r *model.TaskUpdateRequest) error {
	query := `
    UPDATE tasks
    SET
      content = $2,
      description = $3,
      due = $4,
      duration = $5,
      priority = $6
    WHERE id = $1
  `

	_, err := db.Exec(query, r.ID, r.Content, r.Description, r.Due, r.Duration, r.Priority)
	if err != nil {
		return fmt.Errorf("database: unable to update task: %v", err)
	}

	return nil
}

func (db *DB) MoveTask(r *model.TaskMoveRequest) error {
	query := `
    UPDATE tasks
    SET
      project_id = $2,
      parent_id = $3
    WHERE id = $1
  `
	_, err := db.Exec(query, r.ID, r.ProjectID, r.ParentID)
	if err != nil {
		return fmt.Errorf("database: unable to move task: %v", err)
	}

	return nil
}

func (db *DB) ReorderTasks(r *model.TaskReorderRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("database: unable to reorder tasks: %v", err)
	}

	for _, task := range r.Tasks {
		query := `
      UPDATE tasks
      SET
        child_order = $2
      WHERE id = $1
    `
		_, err := tx.Exec(query, task.ID, task.ChildOrder)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("database: unable to reorder tasks: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("database: unable to reorder tasks: %v", err)
	}

	return nil
}

func (db *DB) DeleteTask(r *model.TaskDeleteRequest) error {
	query := `
    DELETE FROM tasks
    WHERE id = $1
  `
	_, err := db.Exec(query, r.ID)
	if err != nil {
		return fmt.Errorf("database: unable to delete task: %v", err)
	}

	return nil
}

func (db *DB) DoneTask(r *model.TaskDoneRequest) error {
	query := `
    UPDATE tasks
    SET
      done = TRUE,
      done_at = NOW()
    WHERE id = $1
  `
	_, err := db.Exec(query, r.ID)
	if err != nil {
		return fmt.Errorf("database: unable to mark task as done: %v", err)
	}

	return nil
}

func (db *DB) UnDoneTask(r *model.TaskUnDoneRequest) error {
	query := `
    UPDATE tasks
    SET
      done = FALSE,
      done_at = NULL
    WHERE id = $1
  `
	_, err := db.Exec(query, r.ID)
	if err != nil {
		return fmt.Errorf("database: unable to mark task as undone: %v", err)
	}

	return nil
}

func (db *DB) ArchiveTask(r *model.TaskDeleteRequest) error {
	query := `
    UPDATE tasks
    SET
      archive = TRUE,
      archive_at = NOW()
    WHERE id = $1
  `
	_, err := db.Exec(query, r.ID)
	if err != nil {
		return fmt.Errorf("database: unable to archive task: %v", err)
	}

	return nil
}

func (db *DB) UnArchiveTask(r *model.TaskDeleteRequest) error {
	query := `
    UPDATE tasks
    SET
      archive = FALSE,
      archive_at = NULL
    WHERE id = $1
  `
	_, err := db.Exec(query, r.ID)
	if err != nil {
		return fmt.Errorf("database: unable to unarchive task: %v", err)
	}

	return nil
}

func (db *DB) GetTask(r *model.TaskGetRequest) (*model.Task, error) {
	query := `
    SELECT
      id,
      user_id,
      content,
      description,
      due,
      duration,
      priority,
      project_id,
      parent_id,
      child_order,
      done,
      done_at,
      archive,
      archive_at,
      created_at,
      updated_at
    FROM tasks
    WHERE id = $1
  `

	t := &model.Task{}
	err := db.QueryRow(query, r.ID).Scan(
		&t.ID,
		&t.UserID,
		&t.Content,
		&t.Description,
		&t.Due,
		&t.Duration,
		&t.Priority,
		&t.ProjectID,
		&t.ParentID,
		&t.ChildOrder,
		&t.Done,
		&t.DoneAt,
		&t.Archive,
		&t.ArchiveAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get task: %v", err)
	}

	return t, nil
}
