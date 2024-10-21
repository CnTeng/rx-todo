SELECT
  id,
  user_id,
  name,
  description,
  child_order,
  inbox,
  favorite,
  total_tasks,
  done_tasks,
  archived,
  archived_at,
  created_at,
  updated_at
FROM
  projects_with_sub_tasks
WHERE
  id = $1
  AND user_id = $2
