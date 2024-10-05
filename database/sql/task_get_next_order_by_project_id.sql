SELECT
  coalesce(max(child_order) + 1, 0)
FROM
  tasks
WHERE
  user_id = $1
  AND (
    project_id = $2
    OR (
      $2 IS NULL
      AND project_id IS NULL
    )
  )
