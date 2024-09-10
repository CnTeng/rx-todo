UPDATE
    projects
SET
    archived = FALSE,
    archived_at = NULL
WHERE
    id = $1
    AND user_id = $2
