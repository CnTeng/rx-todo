INSERT INTO projects (user_id, content, inbox)
    VALUES ($1, $2, $3)
RETURNING
    id
