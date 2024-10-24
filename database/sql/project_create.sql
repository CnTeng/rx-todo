WITH
  new_positions AS (
    SELECT
      coalesce(floor(max("position")) + 1, 0) AS new_position
    FROM
      projects
    WHERE
      user_id = $1
  )
INSERT INTO
  projects (user_id, name, description, favorite, "position")
VALUES
  (
    $1,
    $2,
    $3,
    $4,
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
  inbox,
  archived,
  archived_at,
  created_at,
  updated_at
