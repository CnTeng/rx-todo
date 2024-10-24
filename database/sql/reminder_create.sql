INSERT INTO
  reminders (user_id, task_id, due)
VALUES
  ($1, $2, ROW ($3, $4))
RETURNING
  id,
  created_at,
  updated_at
