INSERT INTO labels (user_id, name, color)
    VALUES ($1, $2, $3)
RETURNING
    id, created_at, updated_at
