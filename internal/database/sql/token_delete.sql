DELETE FROM tokens
WHERE
  id = $1
  AND user_id = $2
