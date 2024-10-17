UPDATE tasks
SET
  name = $3,
  description = $4,
  due = ROW ($5, $6),
  duration = ROW ($7, $8),
  priority = $9,
  project_id = $10,
  parent_id = $11,
  child_order = $12,
  updated_at = now()
WHERE
  id = $1
  AND user_id = $2
RETURNING
  updated_at
