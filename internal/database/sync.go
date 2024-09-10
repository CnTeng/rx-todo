package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/lib/pq"
)

//go:embed sql/sync_create.sql

var createSyncStatusQuery string

func (db *DB) createSyncStatus(tx *sql.Tx, status *model.SyncStatus) error {
	if _, err := tx.Exec(
		createSyncStatusQuery,
		status.UserID,
		pq.Array(status.ObjectIDs),
		status.ObjectType,
		status.Operation); err != nil {
		return fmt.Errorf("failed to create sync status: %w", err)
	}

	return nil
}
