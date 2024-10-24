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

	//go:embed sql/project_get_new_position.sql
	getProjectNewPositionQuery string

	//go:embed sql/project_get_top_position.sql
	getProjectTopPositionQuery string

	//go:embed sql/project_update.sql
	updateProjectQuery string

	//go:embed sql/project_update_position.sql
	updateProjectPositionQuery string

	//go:embed sql/project_update_status.sql
	updateProjectStatusQuery string

	//go:embed sql/project_delete.sql
	deleteProjectQuery string
)

func (db *DB) CreateProject(project *model.Project) (*model.Project, error) {
	if err := db.QueryRow(
		createProjectQuery,
		project.UserID,
		project.Name,
		project.Description,
		project.Favorite,
	).Scan(
		&project.ID,
		&project.Position,
		&project.Inbox,
		&project.Archived,
		&project.ArchivedAt,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
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
		&project.Position,
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
			&project.Position,
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
		project.Position,
		project.Inbox,
		project.Favorite,
	).Scan(&project.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

func (db *DB) UpdateProjectPosition(project *model.Project, previousID int64) (*model.Project, error) {
	return project, db.withTx(func(tx *sql.Tx) error {
		if previousID == 0 {
			if err := db.QueryRow(
				getProjectTopPositionQuery,
				project.UserID,
			).Scan(&project.Position); err != nil {
				return fmt.Errorf("failed to get project top position: %w", err)
			}
		} else {
			if err := db.QueryRow(
				getProjectNewPositionQuery,
				previousID,
				project.UserID,
			).Scan(&project.Position); err != nil {
				return fmt.Errorf("failed to get project new position: %w", err)
			}
		}

		if err := tx.QueryRow(
			updateProjectPositionQuery,
			project.ID,
			project.UserID,
		).Scan(&project.UpdatedAt); err != nil {
			return fmt.Errorf("failed to update project position: %w", err)
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
	).Scan(
		&project.ArchivedAt,
		&project.UpdatedAt,
	); err != nil {
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
