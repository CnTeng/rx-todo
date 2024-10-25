UPDATE tasks
SET
  done = $3,
  done_at = CASE
    WHEN $3 THEN now()
    ELSE NULL
  END
WHERE
  id = $1
  AND user_id = $2
RETURNING
  done_at
