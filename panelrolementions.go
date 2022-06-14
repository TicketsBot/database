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
	"panel_id" int NOT NULL,
	"role_id" int8 NOT NULL,
	FOREIGN KEY("panel_id") REFERENCES panels("panel_id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("panel_id", "role_id")
);
CREATE INDEX IF NOT EXISTS panel_role_mentions_panel_id ON panel_role_mentions("panel_id");
`
}

func (p *PanelRoleMentions) GetRoles(panelId int) (roles []uint64, e error) {
	query := `SELECT "role_id" from panel_role_mentions WHERE "panel_id"=$1;`

	rows, err := p.Query(context.Background(), query, panelId)
	defer rows.Close()
	if err != nil {
		e = err
		return
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			continue
		}

		roles = append(roles, roleId)
	}

	return
}

func (p *PanelRoleMentions) Add(panelId int, roleId uint64) (err error) {
	query := `INSERT INTO panel_role_mentions("panel_id", "role_id") VALUES($1, $2) ON CONFLICT("panel_id", "role_id") DO NOTHING;`
	_, err = p.Exec(context.Background(), query, panelId, roleId)
	return
}

func (p *PanelRoleMentions) DeleteAll(panelId int) (err error) {
	query := `DELETE FROM panel_role_mentions WHERE "panel_id"=$1;`
	_, err = p.Exec(context.Background(), query, panelId)
	return
}

func (p *PanelRoleMentions) DeleteAllRole(roleId uint64) (err error) {
	query := `DELETE FROM panel_role_mentions WHERE "role_id"=$1;`
	_, err = p.Exec(context.Background(), query, roleId)
	return
}

func (p *PanelRoleMentions) Delete(panelId int, roleId uint64) (err error) {
	query := `DELETE FROM panel_role_mentions WHERE "panel_id"=$1 AND "role_id"=$2;`
	_, err = p.Exec(context.Background(), query, panelId, roleId)
	return
}

func (p *PanelRoleMentions) Replace(panelId int, roleIds []uint64) error {
	tx, err := p.Begin(context.Background())
	if err != nil {
		return err
	}

	// Remove existing mentions from panel
	if _, err = tx.Exec(context.Background(), `DELETE FROM panel_role_mentions WHERE "panel_id" = $1;`, panelId); err != nil {
		return err
	}

	// Add each provided mention to panel
	for _, roleId := range roleIds {
		query := `INSERT INTO panel_role_mentions("panel_id", "role_id") VALUES($1, $2) ON CONFLICT("panel_id", "role_id") DO NOTHING;`
		if _, err = tx.Exec(context.Background(), query, panelId, roleId); err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}
