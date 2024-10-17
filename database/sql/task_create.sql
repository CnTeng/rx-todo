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
    child_order
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
    $11
  )
RETURNING
  id,
  created_at,
  updated_at
