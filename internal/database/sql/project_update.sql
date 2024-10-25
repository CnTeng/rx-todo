UPDATE projects
SET
  name = $3,
  description = $4,
  "position" = $5,
  inbox = $6,
  favorite = $7,
  updated_at = now()
WHERE
  id = $1
  AND user_id = $2
RETURNING
  updated_at
