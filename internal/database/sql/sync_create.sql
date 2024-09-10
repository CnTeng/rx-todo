INSERT INTO sync_status (user_id, object_id, object_type, operation)
    VALUES ($1, unnest($2::bigint[]), $3, $4)
