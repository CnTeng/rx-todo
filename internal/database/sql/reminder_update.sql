UPDATE
    reminders
SET
    due = ROW ($3,
        $4),
    updated_at = NOW()
WHERE
    id = $1
    AND user_id = $2
RETURNING
    updated_at
