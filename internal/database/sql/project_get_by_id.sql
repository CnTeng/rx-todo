SELECT
    id,
    user_id,
    content,
    description,
    parent_id,
    child_order,
    inbox,
    favorite,
    archived,
    archived_at,
    created_at,
    updated_at
FROM
    projects
WHERE
    id = $1
    AND user_id = $2
