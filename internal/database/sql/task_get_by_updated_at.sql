SELECT
  id,
  user_id,
  content,
  description,
  (due).date,
  (due).recurring,
  (duration).amount,
  (duration).unit,
  priority,
  project_id,
  parent_id,
  child_order,
  labels,
  done,
  done_at,
  archived,
  archived_at,
  created_at,
  updated_at
FROM
  task_with_labels
WHERE
  user_id = $1
  AND updated_at > $2