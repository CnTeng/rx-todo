UPDATE users
SET
  inbox_id = $2
WHERE
  id = $1
