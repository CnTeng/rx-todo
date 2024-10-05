UPDATE labels
SET
  name = $3,
  color = $4,
  updated_at = now()
WHERE
  id = $1
  AND user_id = $2
RETURNING
  created_at,
  updated_at
