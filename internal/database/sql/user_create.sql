INSERT INTO
  users (username, PASSWORD, email, timezone)
VALUES
  (lower($1), $2, $3, coalesce($4, 'UTC'))
RETURNING
  id,
  timezone,
  created_at,
  updated_at
