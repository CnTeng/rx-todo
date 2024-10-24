SELECT
  id,
  user_id,
  task_id,
  (due).date,
  (due).recurring,
  created_at,
  updated_at
FROM
  reminders
WHERE
  id = $1
  AND user_id = $2
