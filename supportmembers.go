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

func (s *SupportTeamMembersTable) Get(ctx context.Context, teamId int) (members []uint64, e error) {
	rows, err := s.Query(ctx, `SELECT "user_id" from support_team_members WHERE "team_id" = $1;`, teamId)
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

func (s *SupportTeamMembersTable) Add(ctx context.Context, teamId int, userId uint64) (err error) {
	query := `INSERT INTO support_team_members("team_id", "user_id") VALUES($1, $2) ON CONFLICT (team_id, user_id) DO NOTHING;`
	_, err = s.Exec(ctx, query, teamId, userId)
	return
}

func (s *SupportTeamMembersTable) Delete(ctx context.Context, teamId int, userId uint64) (err error) {
	_, err = s.Exec(ctx, `DELETE FROM support_team_members WHERE "team_id"=$1 AND "user_id"=$2;`, teamId, userId)
	return
}

func (s *SupportTeamMembersTable) IsSupport(ctx context.Context, guildId, userId uint64) (isSupport bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team_members
	INNER JOIN support_team
	ON support_team_members.team_id = support_team.id
	WHERE support_team.guild_id = $1 AND support_team_members.user_id = $2
);
`

	err = s.QueryRow(ctx, query, guildId, userId).Scan(&isSupport)
	return
}

func (s *SupportTeamMembersTable) IsSupportSubset(ctx context.Context, guildId, userId uint64, teams []int) (isSupport bool, err error) {
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

	err = s.QueryRow(ctx, query, guildId, userId, teamIdArray).Scan(&isSupport)
	return
}

func (s *SupportTeamMembersTable) GetAllSupportMembers(ctx context.Context, guildId uint64) (users []uint64, err error) {
	query := `
SELECT support_team_members.user_id
FROM support_team_members
INNER JOIN support_team
ON support_team_members.team_id = support_team.id
WHERE support_team.guild_id = $1;
`

	rows, err := s.Query(ctx, query, guildId)
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

func (s *SupportTeamMembersTable) GetAllSupportMembersForPanel(ctx context.Context, panelId int) (users []uint64, err error) {
	query := `
SELECT DISTINCT support_team_members.user_id
FROM support_team_members
INNER JOIN panel_teams
ON support_team_members.team_id = panel_teams.team_id
WHERE panel_teams.panel_id = $1;
`

	rows, err := s.Query(ctx, query, panelId)
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

func (s *SupportTeamMembersTable) GetAllTeamsForUser(ctx context.Context, guildId, userId uint64) ([]int, error) {
	query := `
SELECT support_team_members.team_id
FROM support_team_members
INNER JOIN support_team
ON support_team_members.team_id = support_team.id
WHERE support_team.guild_id = $1 AND support_team_members.user_id = $2;
`

	rows, err := s.Query(ctx, query, guildId, userId)
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
