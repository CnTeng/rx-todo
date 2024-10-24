SELECT
  min("position") / 2
FROM
  projects
WHERE
  user_id = $1
