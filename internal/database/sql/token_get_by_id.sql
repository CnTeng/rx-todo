SELECT
  id,
  user_id,
  token,
  created_at,
  updated_at
FROM
  tokens
WHERE
  user_id = $1
  AND id = $2
