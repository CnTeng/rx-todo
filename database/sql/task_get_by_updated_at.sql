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
  child_order,
  total_tasks,
  done_tasks,
  done,
  done_at,
  archived,
  archived_at,
  created_at,
  updated_at
FROM
  tasks_with_sub_tasks
WHERE
  user_id = $1
  AND updated_at > $2
