SELECT credits
FROM vote_credits
WHERE user_id = $1
FOR UPDATE;