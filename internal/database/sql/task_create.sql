INSERT INTO tasks (user_id, content, description, due, duration, priority, project_id, child_order)
    VALUES ($1, $2, $3, ROW ($4, $5), ROW ($6, $7), $8, $9, $10)
RETURNING
    id, created_at, updated_at
