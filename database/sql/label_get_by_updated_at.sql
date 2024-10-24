SELECT
  id,
  user_id,
  name,
  color,
  created_at,
  updated_at
FROM
  labels
WHERE
  user_id = $1
  AND updated_at > $2
