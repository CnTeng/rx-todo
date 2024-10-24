UPDATE projects
SET
  archived = $3,
  archived_at = CASE
    WHEN $3 THEN now()
    ELSE NULL
  END,
  updated_at = now()
WHERE
  id = $1
  AND user_id = $2
RETURNING
  archived_at,
  updated_at
