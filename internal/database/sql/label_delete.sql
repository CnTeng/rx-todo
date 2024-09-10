DELETE FROM labels
WHERE id = $1
    AND user_id = $2
