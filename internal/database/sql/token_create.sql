INSERT INTO
  tokens (user_id, token)
VALUES
  ($1, $2)
RETURNING
  id,
  token,
  created_at,
  updated_at
