SELECT
  coalesce(max(child_order) + 1)
FROM
  projects
WHERE
  user_id = $1
