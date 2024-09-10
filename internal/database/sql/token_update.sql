UPDATE
    tokens
SET
    token = $3,
    updated_at = NOW()
WHERE
    id = $1
    AND user_id = $2
RETURNING
    id,
    token,
    created_at,
    updated_at
