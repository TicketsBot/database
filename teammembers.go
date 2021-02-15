package database

import (
	"context"
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
	_, err = s.Exec(context.Background(), `DELETE FROM support_team_members WHERE "team_id"=$1 AND "user_id"=$1;`, teamId, userId)
	return
}
