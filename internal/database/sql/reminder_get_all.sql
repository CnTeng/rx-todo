SELECT
    id,
    user_id,
    task_id,
    (due).date,
    (due).recurring,
    created_at,
    updated_at
FROM
    reminders
WHERE
    user_id = $1
