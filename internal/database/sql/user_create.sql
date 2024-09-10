INSERT INTO users (username, PASSWORD, email, timezone)
    VALUES (LOWER($1), $2, $3, COALESCE($4, 'UTC'))
RETURNING
    id, timezone, created_at, updated_at
