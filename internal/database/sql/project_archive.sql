UPDATE
    projects
SET
    archived = TRUE,
    archived_at = NOW()
WHERE
    id = $1
    AND user_id = $2
