UPDATE reminders
SET
  due = ROW ($3, $4),
  updated_at = now()
WHERE
  id = $1
  AND user_id = $2
RETURNING
  updated_at
