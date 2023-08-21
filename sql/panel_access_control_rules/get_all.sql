SELECT role_id, action
FROM panel_access_control_rules
WHERE panel_id = $1
ORDER BY position ASC;