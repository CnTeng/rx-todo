SELECT
  coalesce(max(child_order) + 1)
FROM
  projects
WHERE
  user_id = $1
  AND (
    parent_id = $2
    OR (
      $2 IS NULL
      AND parent_id IS NULL
    )
  )
