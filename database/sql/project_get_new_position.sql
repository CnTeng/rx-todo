WITH
  sorted_projects AS (
    SELECT
      id,
      "position",
      lead("position") OVER (
        ORDER BY
          "position"
      ) AS next_position
    FROM
      projects
    WHERE
      user_id = $2
  )
SELECT
  CASE
    WHEN next_position IS NOT NULL THEN ("position" + next_position) / 2
    ELSE floor("position") + 1
  END AS new_position
FROM
  sorted_projects
WHERE
  id = $1
