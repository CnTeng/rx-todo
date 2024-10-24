INSERT INTO
  task_labels (task_id, label_id)
SELECT
  $3,
  id
FROM
  labels
WHERE
  user_id = $1
  AND name = $2
