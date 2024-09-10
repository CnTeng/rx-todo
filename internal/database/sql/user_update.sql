UPDATE
    users
SET
    username = LOWER($2),
    PASSWORD = $3,
    email = $4,
    timezone = $5,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    created_at,
    updated_at
