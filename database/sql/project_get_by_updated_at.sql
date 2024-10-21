SELECT
  id,
  user_id,
  name,
  description,
  child_order,
  inbox,
  favorite,
  archived,
  archived_at,
  created_at,
  updated_at
FROM
  projects
WHERE
  user_id = $1
  AND updated_at > $2
