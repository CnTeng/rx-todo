WITH
  new_positions AS (
    SELECT
      coalesce(floor(max("position")) + 1, 0) AS new_position
    FROM
      tasks
    WHERE
      user_id = $1
      AND (
        (
          $9 IS NOT NULL
          AND project_id = $9
        )
        OR (
          $10 IS NOT NULL
          AND parent_id = $10
        )
      )
  )
INSERT INTO
  tasks (
    user_id,
    name,
    description,
    due,
    duration,
    priority,
    project_id,
    parent_id,
    "position"
  )
VALUES
  (
    $1,
    $2,
    $3,
    ROW ($4, $5),
    ROW ($6, $7),
    $8,
    $9,
    $10,
    (
      SELECT
        new_position
      FROM
        new_positions
    )
  )
RETURNING
  id,
  "position",
  created_at,
  updated_at
