SELECT panel_access_control_rules.panel_id,
       panel_access_control_rules.role_id,
       panel_access_control_rules.action
FROM panel_access_control_rules
         INNER JOIN panels ON panels.panel_id = panel_access_control_rules.panel_id
WHERE panels.guild_id = $1
ORDER BY panel_access_control_rules.position ASC;