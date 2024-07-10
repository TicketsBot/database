package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SupportTeamRolesTable struct {
	*pgxpool.Pool
}

func newSupportTeamRolesTable(db *pgxpool.Pool) *SupportTeamRolesTable {
	return &SupportTeamRolesTable{
		db,
	}
}

func (s SupportTeamRolesTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS support_team_roles(
	"team_id" int NOT NULL,
	"role_id" int8 NOT NULL,
	FOREIGN KEY("team_id") REFERENCES support_team("id") ON DELETE CASCADE ON UPDATE CASCADE,
	PRIMARY KEY("team_id", "role_id")
);`
}

func (s *SupportTeamRolesTable) Get(ctx context.Context, teamId int) (roles []uint64, e error) {
	rows, err := s.Query(ctx, `SELECT "role_id" from support_team_roles WHERE "team_id" = $1;`, teamId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			return nil, err
		}

		roles = append(roles, roleId)
	}

	return
}

func (s *SupportTeamRolesTable) Add(ctx context.Context, teamId int, roleId uint64) (err error) {
	query := `INSERT INTO support_team_roles("team_id", "role_id") VALUES($1, $2) ON CONFLICT (team_id, role_id) DO NOTHING;`
	_, err = s.Exec(ctx, query, teamId, roleId)
	return
}

func (s *SupportTeamRolesTable) Delete(ctx context.Context, teamId int, roleId uint64) (err error) {
	_, err = s.Exec(ctx, `DELETE FROM support_team_roles WHERE "team_id"=$1 AND "role_id"=$2;`, teamId, roleId)
	return
}

func (s *SupportTeamRolesTable) DeleteAllRole(ctx context.Context, roleId uint64) (err error) {
	_, err = s.Exec(ctx, `DELETE FROM support_team_roles WHERE "role_id"=$1;`, roleId)
	return
}

func (s *SupportTeamRolesTable) IsSupport(ctx context.Context, guildId, roleId uint64) (isSupport bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team_roles
	INNER JOIN support_team
	ON support_team_roles.team_id = support_team.id
	WHERE support_team.guild_id = $1 AND support_team_roles.role_id = $2
);
`

	err = s.QueryRow(ctx, query, guildId, roleId).Scan(&isSupport)
	return
}

func (s *SupportTeamRolesTable) IsSupportAny(ctx context.Context, guildId uint64, roleIds []uint64) (isSupport bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team_roles
	INNER JOIN support_team
	ON support_team_roles.team_id = support_team.id
	WHERE support_team.guild_id = $1 AND support_team_roles.role_id = ANY($2)
);
`

	roleIdArray := &pgtype.Int8Array{}
	if err = roleIdArray.Set(roleIds); err != nil {
		return
	}

	err = s.QueryRow(ctx, query, guildId, roleIdArray).Scan(&isSupport)
	return
}

func (s *SupportTeamRolesTable) IsSupportAnySubset(ctx context.Context, guildId uint64, roleIds []uint64, teamIds []int) (isSupport bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team_roles
	INNER JOIN support_team
	ON support_team_roles.team_id = support_team.id
	WHERE support_team.guild_id = $1 AND support_team_roles.role_id = ANY($2) AND support_team.id = ANY($3)
);
`

	roleIdArray := &pgtype.Int8Array{}
	if err := roleIdArray.Set(roleIds); err != nil {
		return false, err
	}

	teamIdArray := &pgtype.Int4Array{}
	if err := teamIdArray.Set(teamIds); err != nil {
		return false, err
	}

	err = s.QueryRow(ctx, query, guildId, roleIdArray, teamIdArray).Scan(&isSupport)
	return
}

func (s *SupportTeamRolesTable) GetAllSupportRoles(ctx context.Context, guildId uint64) (roles []uint64, err error) {
	query := `
SELECT support_team_roles.role_id
FROM support_team_roles
INNER JOIN support_team
ON support_team_roles.team_id = support_team.id
WHERE support_team.guild_id = $1;
`

	rows, err := s.Query(ctx, query, guildId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			return nil, err
		}

		roles = append(roles, roleId)
	}

	return
}

func (s *SupportTeamRolesTable) GetAllSupportRolesForPanel(ctx context.Context, panelId int) (roles []uint64, err error) {
	query := `
SELECT DISTINCT support_team_roles.role_id
FROM support_team_roles
INNER JOIN panel_teams
ON support_team_roles.team_id = panel_teams.team_id
WHERE panel_teams.panel_id = $1;
`

	rows, err := s.Query(ctx, query, panelId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			return nil, err
		}

		roles = append(roles, roleId)
	}

	return
}

func (s *SupportTeamRolesTable) GetAllTeamsForRoles(ctx context.Context, guildId uint64, roleIds []uint64) ([]int, error) {
	query := `
SELECT support_team_roles.team_id
FROM support_team_roles
INNER JOIN support_team
ON support_team_roles.team_id = support_team.id
WHERE support_team.guild_id = $1 AND support_team_roles.role_id = ANY($2);
`

	roleIdArray := &pgtype.Int8Array{}
	if err := roleIdArray.Set(roleIds); err != nil {
		return nil, err
	}

	rows, err := s.Query(ctx, query, guildId, roleIdArray)
	if err != nil {
		return nil, err
	}

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
