INSERT INTO
  projects (
    user_id,
    content,
    description,
    parent_id,
    child_order,
    favorite
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  id,
  inbox,
  archived,
  archived_at,
  created_at,
  updated_at
