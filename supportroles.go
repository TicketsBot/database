package database

import (
	"context"
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
