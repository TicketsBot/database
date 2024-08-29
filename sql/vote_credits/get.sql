SELECT COALESCE(credits, 0)
FROM vote_credits
WHERE user_id = $1
FOR UPDATE;