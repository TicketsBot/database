CREATE TABLE IF NOT EXISTS panel_access_control_rules
(
    "panel_id" int        NOT NULL,
    "role_id"  int8       NOT NULL,
    "position" int        NOT NULL,
    "action"   varchar(5) NOT NULL,
    UNIQUE ("panel_id", "role_id"),
    UNIQUE ("panel_id", "position"),
    FOREIGN KEY ("panel_id") REFERENCES panels ("panel_id") ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY ("panel_id", "role_id")
);
CREATE INDEX IF NOT EXISTS panel_access_control_rules_panel_id ON panel_access_control_rules ("panel_id");
