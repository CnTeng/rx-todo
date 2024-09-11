UPDATE projects
SET
  archived = TRUE,
  archived_at = now()
WHERE
  id = $1
  AND user_id = $2
