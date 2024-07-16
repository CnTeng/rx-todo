package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) CreateProject(project *model.Project) (*model.Project, error) {
	query := `
		INSERT INTO projects
			(user_id, content, description, parent_id, child_order, favorite)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		RETURNING 
      id, inbox, archived, archived_at, created_at, updated_at
	`

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("database: unable to begin transaction: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	if project.ParentID == nil {
		_, err = tx.Exec(`
			UPDATE 
				projects 
			SET 
				child_order = child_order + 1
    	WHERE 
				user_id = $1 AND parent_id IS NULL AND child_order >= $2
			`, project.UserID, project.ChildOrder)
	} else {
		_, err = tx.Exec(`
			UPDATE 
				projects 
			SET 
				child_order = child_order + 1
    	WHERE 
				user_id = $1 AND parent_id = $2 AND child_order >= $3
			`, project.UserID, project.ParentID, project.ChildOrder)
	}
	if err != nil {
		return nil, fmt.Errorf("database: unable to create task: %v", err)
	}

	err = tx.QueryRow(
		query,
		project.UserID,
		project.Content,
		project.Description,
		project.ParentID,
		project.ChildOrder,
		project.Favorite,
	).Scan(
		&project.ID,
		&project.Inbox,
		&project.Archived,
		&project.ArchivedAt,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create task: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("database: unable to commit transaction: %v", err)
	}

	return project, nil
}

func (db *DB) GetProjectByID(userID, id int64) (*model.Project, error) {
	project := &model.Project{}
	query := `
		SELECT
      id,
      user_id,
      content,
      description,
      parent_id,
      child_order,
      inbox,
      favorite,
      archived,
      archived_at,
      created_at,
      updated_at
		FROM
      projects
		WHERE
			user_id = $1 AND id = $2
	`

	err := db.QueryRow(query, userID, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Content,
		&project.Description,
		&project.ParentID,
		&project.ChildOrder,
		&project.Inbox,
		&project.Favorite,
		&project.Archived,
		&project.ArchivedAt,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get task: %v", err)
	}

	return project, nil
}

func (db *DB) GetProjects(userID int64) ([]*model.Project, error) {
	var projects []*model.Project
	query := `
		SELECT
      id,
      user_id,
      content,
      description,
      parent_id,
      child_order,
      inbox,
      favorite,
      archived,
      archived_at,
      created_at,
      updated_at
		FROM
			projects
		WHERE
			user_id = $1
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		project := &model.Project{}

		if err := rows.Scan(

			&project.ID,
			&project.UserID,
			&project.Content,
			&project.Description,
			&project.ParentID,
			&project.ChildOrder,
			&project.Inbox,
			&project.Favorite,
			&project.Archived,
			&project.ArchivedAt,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("database: unable to get tasks: %v", err)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (db *DB) UpdateProject(project *model.Project) (*model.Project, error) {
	query := `
		UPDATE 
			projects
		SET 
			content = $3,
			description = $4,
			parent_id = $5,
			child_order = $6,
			inbox = $7,
			favorite = $8,
			updated_at = NOW()
		WHERE 
			id = $1 AND user_id = $2
		RETURNING 
			updated_at
	`

	err := db.QueryRow(
		query,
		project.ID,
		project.UserID,
		project.Content,
		project.Description,
		project.ParentID,
		project.ChildOrder,
		project.Inbox,
		project.Favorite,
	).Scan(&project.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to update task: %v", err)
	}

	return project, nil
}

func (db *DB) UpdateProjects(projects []*model.Project) error {
	query := `
		UPDATE 
			projects
		SET 
			content = $3,
			description = $4,
			parent_id = $5,
			child_order = $6,
			updated_at = NOW()
		WHERE 
			id = $1 AND user_id = $2
		RETURNING 
			updated_at
	`

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("database: unable to begin transaction: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, project := range projects {
		err := tx.QueryRow(
			query,
			project.ID,
			project.UserID,
			project.Content,
			project.Description,
			project.ParentID,
			project.ChildOrder,
		).Scan(&project.UpdatedAt)
		if err != nil {
			return fmt.Errorf("database: unable to update tasks: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("database: unable to commit transaction: %v", err)
	}

	return nil
}

func (db *DB) DeleteProject(userID, id int64) error {
	query := `DELETE FROM projects WHERE id = $1 AND user_id = $2`
	return db.execSimpleQuery(query, id, userID)
}

func (db *DB) ArchiveProject(userID, id int64) error {
	query := `
		UPDATE 
			projects 
		SET 
			archived = TRUE,
			archived_at = NOW()
		WHERE
			id = $1 AND user_id = $2
	`
	return db.execSimpleQuery(query, id, userID)
}

func (db *DB) UnarchiveProject(userID, id int64) error {
	query := `
		UPDATE 
			projects 
		SET 
			archived = FALSE,
			archived_at = NULL
		WHERE
			id = $1 AND user_id = $2
	`
	return db.execSimpleQuery(query, id, userID)
}
