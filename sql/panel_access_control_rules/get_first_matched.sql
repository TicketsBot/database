SELECT role_id, action
FROM panel_access_control_rules
WHERE panel_id = $1 AND role_id = ANY($2)
ORDER BY position ASC
LIMIT 1;