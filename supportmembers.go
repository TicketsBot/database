package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SupportTeamMembersTable struct {
	*pgxpool.Pool
}

func newSupportTeamMembersTable(db *pgxpool.Pool) *SupportTeamMembersTable {
	return &SupportTeamMembersTable{
		db,
	}
}

func (s SupportTeamMembersTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS support_team_members(
	"team_id" int NOT NULL,
	"user_id" int8 NOT NULL,
	FOREIGN KEY("team_id") REFERENCES support_team("id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("team_id", "user_id")
);`
}

func (s *SupportTeamMembersTable) Get(teamId int) (members []uint64, e error) {
	rows, err := s.Query(context.Background(), `SELECT "user_id" from support_team_members WHERE "team_id" = $1;`, teamId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}

		members = append(members, userId)
	}

	return
}

func (s *SupportTeamMembersTable) Add(teamId int, userId uint64) (err error) {
	_, err = s.Exec(context.Background(), `INSERT INTO support_team_members("team_id", "user_id") VALUES($1, $2);`, teamId, userId)
	return
}

func (s *SupportTeamMembersTable) Delete(teamId int, userId uint64) (err error) {
	_, err = s.Exec(context.Background(), `DELETE FROM support_team_members WHERE "team_id"=$1 AND "user_id"=$2;`, teamId, userId)
	return
}

func (s *SupportTeamMembersTable) IsSupport(guildId, userId uint64) (isSupport bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team_members
	INNER JOIN support_team
	ON support_team_members.team_id = support_team.id
	WHERE support_team.guild_id = $1 AND support_team_members.user_id = $2
);
`

	err = s.QueryRow(context.Background(), query, guildId, userId).Scan(&isSupport)
	return
}

func (s *SupportTeamMembersTable) IsSupportSubset(guildId, userId uint64, teams []int) (isSupport bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team_members
	INNER JOIN support_team
	ON support_team_members.team_id = support_team.id
	WHERE support_team.guild_id = $1 AND support_team_members.user_id = $2 AND support_team.id = ANY($3)
);
`

	teamIdArray := &pgtype.Int4Array{}
	if err := teamIdArray.Set(teams); err != nil {
		return false, err
	}

	err = s.QueryRow(context.Background(), query, guildId, userId, teamIdArray).Scan(&isSupport)
	return
}

func (s *SupportTeamMembersTable) GetAllSupportMembers(guildId uint64) (users []uint64, err error) {
	query := `
SELECT support_team_members.user_id
FROM support_team_members
INNER JOIN support_team
ON support_team_members.team_id = support_team.id
WHERE support_team.guild_id = $1;
`

	rows, err := s.Query(context.Background(), query, guildId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}

		users = append(users, userId)
	}

	return
}

func (s *SupportTeamMembersTable) GetAllSupportMembersForPanel(panelId int) (users []uint64, err error) {
	query := `
SELECT DISTINCT support_team_members.user_id
FROM support_team_members
INNER JOIN panel_teams
ON support_team_members.team_id = panel_teams.team_id
WHERE panel_teams.panel_id = $1;
`

	rows, err := s.Query(context.Background(), query, panelId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}

		users = append(users, userId)
	}

	return
}

func (s *SupportTeamMembersTable) GetAllTeamsForUser(guildId, userId uint64) ([]int, error) {
	query := `
SELECT support_team_members.team_id
FROM support_team_members
INNER JOIN support_team
ON support_team_members.team_id = support_team.id
WHERE support_team.guild_id = $1 AND support_team_members.user_id = $2;
`

	rows, err := s.Query(context.Background(), query, guildId, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var teamIds []int
	for rows.Next() {
		var teamId int
		if err := rows.Scan(&teamId); err != nil {
			return nil, err
		}

		teamIds = append(teamIds, teamId)
	}

	return teamIds, nil
}
