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

func (s *SupportTeamRolesTable) Get(teamId int) (roles []uint64, e error) {
	rows, err := s.Query(context.Background(), `SELECT "role_id" from support_team_roles WHERE "team_id" = $1;`, teamId)
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

func (s *SupportTeamRolesTable) Add(teamId int, roleId uint64) (err error) {
	_, err = s.Exec(context.Background(), `INSERT INTO support_team_roles("team_id", "role_id") VALUES($1, $2);`, teamId, roleId)
	return
}

func (s *SupportTeamRolesTable) Delete(teamId int, roleId uint64) (err error) {
	_, err = s.Exec(context.Background(), `DELETE FROM support_team_roles WHERE "team_id"=$1 AND "role_id"=$2;`, teamId, roleId)
	return
}

func (s *SupportTeamRolesTable) DeleteAllRole(roleId uint64) (err error) {
	_, err = s.Exec(context.Background(), `DELETE FROM support_team_roles WHERE "role_id"=$1;`, roleId)
	return
}

func (s *SupportTeamRolesTable) IsSupport(guildId, roleId uint64) (isSupport bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team_roles
	INNER JOIN support_team
	ON support_team_roles.team_id = support_team.id
	WHERE support_team.guild_id = $1 AND support_team_roles.role_id = $2
);
`

	err = s.QueryRow(context.Background(), query, guildId, roleId).Scan(&isSupport)
	return
}

func (s *SupportTeamRolesTable) IsSupportAny(guildId uint64, roleIds []uint64) (isSupport bool, err error) {
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

	err = s.QueryRow(context.Background(), query, guildId, roleIdArray).Scan(&isSupport)
	return
}

func (s *SupportTeamRolesTable) IsSupportAnySubset(guildId uint64, roleIds []uint64, teamIds []int) (isSupport bool, err error) {
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

	err = s.QueryRow(context.Background(), query, guildId, roleIdArray, teamIdArray).Scan(&isSupport)
	return
}

func (s *SupportTeamRolesTable) GetAllSupportRoles(guildId uint64) (roles []uint64, err error) {
	query := `
SELECT support_team_roles.role_id
FROM support_team_roles
INNER JOIN support_team
ON support_team_roles.team_id = support_team.id
WHERE support_team.guild_id = $1;
`

	rows, err := s.Query(context.Background(), query, guildId)
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