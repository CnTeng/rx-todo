package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

var (
	//go:embed sql/project_create.sql
	createProjectQuery string

	//go:embed sql/project_get_by_id.sql
	getProjectByIDQuery string

	//go:embed sql/project_get_all.sql
	getProjectsQuery string

	//go:embed sql/project_update.sql
	updateProjectQuery string

	//go:embed sql/project_update_order.sql
	updateProjectOrderQuery string

	//go:embed sql/project_delete.sql
	deleteProjectQuery string

	//go:embed sql/project_archive.sql
	archiveProjectQuery string

	//go:embed sql/project_unarchive.sql
	unarchiveProjectQuery string
)

func (db *DB) CreateProject(project *model.Project) (*model.Project, error) {
	return project, db.withTx(
		func(tx *sql.Tx) (*model.SyncStatus, error) {
			rows, err := tx.Query(
				updateProjectOrderQuery,
				project.UserID,
				project.ParentID,
				project.ChildOrder,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update project order: %w", err)
			}
			defer rows.Close()

			status := project.ToSyncStatus(model.UpdateOperation).ClearObjectIDs()
			for rows.Next() {
				var id int64
				if err := rows.Scan(&id); err != nil {
					return nil, fmt.Errorf("failed to scan project id: %w", err)
				}
				status.ObjectIDs = append(status.ObjectIDs, id)
			}

			return status, nil
		},
		func(tx *sql.Tx) (*model.SyncStatus, error) {
			err := tx.QueryRow(
				createProjectQuery,
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

			return project.ToSyncStatus(model.CreateOperation), err
		})
}

func (db *DB) GetProjectByID(userID, id int64) (*model.Project, error) {
	project := &model.Project{}

	err := db.QueryRow(getProjectByIDQuery, id, userID).Scan(
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
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (db *DB) GetProjects(userID int64) ([]*model.Project, error) {
	var projects []*model.Project

	rows, err := db.Query(getProjectsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
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
			return nil, fmt.Errorf("failed to get tasks: %w", err)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (db *DB) UpdateProject(project *model.Project) (*model.Project, error) {
	return project, db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		err := tx.QueryRow(
			updateProjectQuery,
			project.ID,
			project.UserID,
			project.Content,
			project.Description,
			project.ParentID,
			project.ChildOrder,
			project.Inbox,
			project.Favorite,
		).Scan(&project.UpdatedAt)

		return project.ToSyncStatus(model.UpdateOperation), err
	})
}

func (db *DB) UpdateProjects(projects []*model.Project) error {
	return db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		status := projects[0].ToSyncStatus(model.UpdateOperation).ClearObjectIDs()
		for _, project := range projects {
			err := tx.QueryRow(
				updateProjectQuery,
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
				return nil, fmt.Errorf("failed to update projects: %w", err)
			}

			status.ObjectIDs = append(status.ObjectIDs, project.ID)
		}

		return status, nil
	})
}

func (db *DB) DeleteProject(project *model.Project) error {
	if inboxID, err := db.GetUserInboxID(project.UserID); err == nil && inboxID == project.ID {
		return fmt.Errorf("failed to delete inbox")
	}

	return db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		if _, err := db.Exec(deleteProjectQuery, project.ID, project.UserID); err != nil {
			return nil, fmt.Errorf("failed to delete project: %w", err)
		}

		return project.ToSyncStatus(model.DeleteOperation), nil
	})
}

func (db *DB) ArchiveProject(project *model.Project) error {
	return db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		if _, err := db.Exec(archiveProjectQuery, project.ID, project.UserID); err != nil {
			return nil, fmt.Errorf("failed to archive project: %w", err)
		}

		return project.ToSyncStatus(model.UpdateOperation), nil
	})
}

func (db *DB) UnarchiveProject(project *model.Project) error {
	return db.withTx(func(tx *sql.Tx) (*model.SyncStatus, error) {
		if _, err := db.Exec(unarchiveProjectQuery, project.ID, project.UserID); err != nil {
			return nil, fmt.Errorf("failed to unarchive project: %w", err)
		}

		return project.ToSyncStatus(model.UpdateOperation), nil
	})
}
