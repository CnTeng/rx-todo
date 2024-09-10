SELECT
    id,
    user_id,
    name,
    color,
    created_at,
    updated_at
FROM
    labels
WHERE
    id = $1
    AND user_id = $2
