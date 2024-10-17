UPDATE projects
SET
  name = $3,
  description = $4,
  parent_id = $5,
  child_order = $6,
  inbox = $7,
  favorite = $8,
  updated_at = now()
WHERE
  id = $1
  AND user_id = $2
RETURNING
  updated_at
