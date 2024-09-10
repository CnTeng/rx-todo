SELECT
    id,
    username,
    PASSWORD,
    email,
    timezone,
    created_at,
    updated_at
FROM
    users
WHERE
    email = $1
