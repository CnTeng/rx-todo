DELETE FROM reminders
WHERE
  id = $1
  AND user_id = $2
