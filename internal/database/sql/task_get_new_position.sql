WITH
  sorted_tasks AS (
    SELECT
      id,
      "position",
      lead("position") OVER (
        ORDER BY
          "position"
      ) AS next_position
    FROM
      tasks
    WHERE
      user_id = $2,
      (
        $3 IS NOT NULL
        AND project_id = $3
      )
      OR (
        $4 IS NOT NULL
        AND parent_id = $4
      )
  )
SELECT
  CASE
    WHEN next_position IS NOT NULL THEN ("position" + next_position) / 2
    ELSE floor("position") + 1
  END AS new_position
FROM
  sorted_tasks
WHERE
  id = $1
