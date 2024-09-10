UPDATE
    projects
SET
    child_order = child_order + 1
WHERE
    user_id = $1
    AND (parent_id = $2
        OR (parent_id IS NULL
            AND $2 IS NULL))
    AND child_order >= $3
RETURNING
    id
