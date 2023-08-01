INSERT INTO s (user_id, token_hash) 
VALUES ($1, $2)
RETURNING id;