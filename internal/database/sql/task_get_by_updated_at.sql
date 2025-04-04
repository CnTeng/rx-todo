SELECT
  id,
  user_id,
  name,
  description,
  (due).date,
  (due).recurring,
  (duration).amount,
  (duration).unit,
  priority,
  project_id,
  parent_id,
  "position",
  done,
  done_at,
  archived,
  archived_at,
  created_at,
  updated_at
FROM
  tasks
WHERE
  user_id = $1
  AND updated_at > $2
