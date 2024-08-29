INSERT INTO vote_credits (user_id, credits)
VALUES ($1, 1)
ON CONFLICT (user_id) DO UPDATE SET credits = vote_credits.credits + 1;