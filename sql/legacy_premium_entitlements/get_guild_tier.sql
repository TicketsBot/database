SELECT max(ent.tier)
FROM legacy_premium_entitlements as ent
LEFT OUTER JOIN permissions ON permissions.user_id = ent.user_id AND permissions.guild_id = $1
WHERE
    ent.expires_at > (NOW() - $3::interval)
    AND
    (
        ent.user_id = $2
        OR
        (ent.user_id = permissions.user_id AND permissions.admin = 't' AND permissions.guild_id = $1)
    )
GROUP BY ent.user_id
;
