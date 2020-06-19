package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PanelRoleMentions struct {
	*pgxpool.Pool
}

func newPanelRoleMentions(db *pgxpool.Pool) *PanelRoleMentions {
	return &PanelRoleMentions{
		db,
	}
}

func (p PanelRoleMentions) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS panel_role_mentions(
	"panel_message_id" int8 NOT NULL,
	"role_id" int8 NOT NULL,
	FOREIGN KEY("panel_message_id") REFERENCES panels("message_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("panel_message_id", "role_id")
);
CREATE INDEX IF NOT EXISTS panel_role_mentions_panel_message_id ON panel_mentions("panel_message_id");
`
}

func (p *PanelRoleMentions) GetRoles(panelMessageId uint64) (roles []uint64, e error) {
	query := `SELECT "role_id" from panel_role_mentions WHERE "panel_message_id"=$1;`

	rows, err := p.Query(context.Background(), query, panelMessageId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var roleId uint16
		if err := rows.Scan(&roleId); err != nil {
			continue
		}

		roles = append(role, roleId)
	}

	return
}

func (p *PanelRoleMentions) Add(panelMessageId, roleId uint64) (err error) {
	query := `INSERT INTO panel_role_mentions("panel_message_id", "role_id") VALUES($1, $2) ON CONFLICT("panel_message_id", "role_id") DO NOTHING;`
	_, err = p.Exec(context.Background(), query, panelMessageId, roleId)
	return
}

func (p *PanelRoleMentions) Delete(panelMessageId, roleId uint64) (err error) {
	query := `DELETE FROM panel_mentions WHERE "panel_message_id"=$1 AND "role_id"=$2;`
	_, err = p.Exec(context.Background(), query, panelMessageId, roleId)
	return
}
