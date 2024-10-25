SELECT
  min("position") / 2
FROM
  tasks
WHERE
  user_id = $1
