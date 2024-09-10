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
    user_id = $1
