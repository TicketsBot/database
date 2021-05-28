package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PanelTeamsTable struct {
	*pgxpool.Pool
}

func newPanelTeamsTable(db *pgxpool.Pool) *PanelTeamsTable {
	return &PanelTeamsTable{
		db,
	}
}

func (p PanelTeamsTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS panel_teams(
	"panel_id" int NOT NULL,
	"team_id" int NOT NULL,
	FOREIGN KEY("panel_id") REFERENCES panels("panel_id") ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY("team_id") REFERENCES support_team("id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("panel_id", "team_id")
);
CREATE INDEX IF NOT EXISTS panel_teams_panel_id ON panel_teams("panel_id");
`
}

func (p *PanelTeamsTable) GetTeams(panelId int) (teams []SupportTeam, e error) {
	query := `
SELECT support_team.id, support_team.guild_id, support_team.name
FROM panel_teams
INNER JOIN support_team
ON panel_teams.team_id = support_team.id
WHERE panel_teams.panel_id = $1;
`

	rows, err := p.Query(context.Background(), query, panelId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var team SupportTeam
		if err := rows.Scan(&team.Id, &team.GuildId, &team.Name); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return
}

func (p *PanelTeamsTable) Add(panelId, teamId int) (err error) {
	query := `INSERT INTO panel_teams("panel_id", "team_id") VALUES($1, $2) ON CONFLICT("panel_id", "team_id") DO NOTHING;`
	_, err = p.Exec(context.Background(), query, panelId, teamId)
	return
}

func (p *PanelTeamsTable) DeleteAll(panelId int) (err error) {
	query := `DELETE FROM panel_teams WHERE "panel_id"=$1;`
	_, err = p.Exec(context.Background(), query, panelId)
	return
}

func (p *PanelTeamsTable) Delete(panelId int, teamId uint64) (err error) {
	query := `DELETE FROM panel_teams WHERE "panel_id"=$1 AND "team_id"=$2;`
	_, err = p.Exec(context.Background(), query, panelId, teamId)
	return
}
