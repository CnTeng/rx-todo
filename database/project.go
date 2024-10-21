package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/CnTeng/rx-todo/model"
)

var (
	//go:embed sql/project_create.sql
	createProjectQuery string

	//go:embed sql/project_get_by_id.sql
	getProjectByIDQuery string

	//go:embed sql/project_get_all.sql
	getProjectsQuery string

	//go:embed sql/project_get_by_updated_at.sql
	getProjectsByUpdateTimeQuery string

	//go:embed sql/project_get_next_order.sql
	getProjectNextOrderQuery string

	//go:embed sql/project_update.sql
	updateProjectQuery string

	//go:embed sql/project_update_status.sql
	updateProjectStatusQuery string

	//go:embed sql/project_delete.sql
	deleteProjectQuery string
)

func (db *DB) CreateProject(project *model.Project) (*model.Project, error) {
	return project, db.withTx(func(tx *sql.Tx) error {
		if err := tx.QueryRow(
			getProjectNextOrderQuery,
			project.UserID,
		).Scan(&project.ChildOrder); err != nil {
			return fmt.Errorf("failed to get project child_order: %w", err)
		}

		if err := tx.QueryRow(
			createProjectQuery,
			project.UserID,
			project.Name,
			project.Description,
			project.ChildOrder,
			project.Favorite,
		).Scan(
			&project.ID,
			&project.Inbox,
			&project.Archived,
			&project.ArchivedAt,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		return nil
	})
}

func (db *DB) GetProjectByID(id, userID int64) (*model.Project, error) {
	project := &model.Project{}

	if err := db.QueryRow(
		getProjectByIDQuery,
		id,
		userID,
	).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.ChildOrder,
		&project.Inbox,
		&project.Favorite,
		&project.SubTasks.Total,
		&project.SubTasks.Done,
		&project.Archived,
		&project.ArchivedAt,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (db *DB) getProjects(query string, args ...any) ([]*model.Project, error) {
	var projects []*model.Project

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		project := &model.Project{}

		if err := rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&project.Description,
			&project.ChildOrder,
			&project.Inbox,
			&project.Favorite,
			&project.SubTasks.Total,
			&project.SubTasks.Done,
			&project.Archived,
			&project.ArchivedAt,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get tasks: %w", err)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (db *DB) GetProjects(userID int64) ([]*model.Project, error) {
	return db.getProjects(getProjectsQuery, userID)
}

func (db *DB) GetProjectsByUpdateTime(userID int64, updateTime *time.Time) ([]*model.Project, error) {
	return db.getProjects(getProjectsByUpdateTimeQuery, userID, updateTime)
}

func (db *DB) UpdateProject(project *model.Project) (*model.Project, error) {
	if err := db.QueryRow(
		updateProjectQuery,
		project.ID,
		project.UserID,
		project.Name,
		project.Description,
		project.ChildOrder,
		project.Inbox,
		project.Favorite,
	).Scan(&project.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

func (db *DB) UpdateProjects(projects []*model.Project) ([]*model.Project, error) {
	return projects, db.withTx(func(tx *sql.Tx) error {
		for _, project := range projects {
			err := tx.QueryRow(
				updateProjectQuery,
				project.ID,
				project.UserID,
				project.Name,
				project.Description,
				project.ChildOrder,
				project.Inbox,
				project.Favorite,
			).Scan(&project.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to update projects: %w", err)
			}
		}

		return nil
	})
}

func (db *DB) UpdateProjectStatus(project *model.Project) (*model.Project, error) {
	if err := db.QueryRow(
		updateProjectStatusQuery,
		project.ID,
		project.UserID,
		project.Archived,
	).Scan(&project.ArchivedAt); err != nil {
		return nil, fmt.Errorf("failed to update project status: %w", err)
	}
	return project, nil
}

func (db *DB) DeleteProject(id, userID int64) error {
	inboxID, err := db.GetUserInboxID(userID)
	if err == nil && inboxID == id {
		return fmt.Errorf("failed to delete inbox")
	}

	return db.withTx(func(tx *sql.Tx) error {
		if _, err := tx.Exec(deleteProjectQuery, id, userID); err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}

		if _, err := tx.Exec(createDeletionLogQuery, userID, "project", id); err != nil {
			return fmt.Errorf("failed to create deletion log: %w", err)
		}
		return nil
	})
}
