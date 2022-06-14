package database

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SupportTeamTable struct {
	*pgxpool.Pool
}

type SupportTeam struct {
	Id      int    `json:"id"`
	GuildId uint64 `json:"guild_id"`
	Name    string `json:"name"`
}

func NewSupportTeam(id int, guildId uint64, name string) SupportTeam {
	return SupportTeam{
		Id:      id,
		GuildId: guildId,
		Name:    name,
	}
}

func newSupportTeamTable(db *pgxpool.Pool) *SupportTeamTable {
	return &SupportTeamTable{
		db,
	}
}

func (s SupportTeamTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS support_team(
	"id" SERIAL NOT NULL UNIQUE,
	"guild_id" int8 NOT NULL,
	"name" VARCHAR(32) NOT NULL,
	UNIQUE("guild_id", "name"),
	PRIMARY KEY("id")
);`
}

func (s *SupportTeamTable) Exists(teamId int, guildId uint64) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM support_team WHERE "id" = $1 and "guild_id" = $2);`
	err = s.QueryRow(context.Background(), query, teamId, guildId).Scan(&exists)
	return
}

func (s *SupportTeamTable) AllTeamsMatchGuild(guildId uint64, teams []int) (valid bool, err error) {
	query := `SELECT NOT EXISTS(SELECT 1 FROM support_team WHERE "id" = ANY($1) and "guild_id" != $2);`

	array := &pgtype.Int4Array{}
	if err := array.Set(teams); err != nil {
		return false, err
	}

	err = s.QueryRow(context.Background(), query, array, guildId).Scan(&valid)
	return
}

func (s *SupportTeamTable) AllTeamsExistForGuild(guildId uint64, teams []int) (valid bool, err error) {
	query := `
SELECT EXISTS(
	SELECT 1
	FROM support_team
	WHERE "guild_id" = $1
	GROUP BY "guild_id"
	HAVING array_agg(id) @> $2
);
`

	array := &pgtype.Int4Array{}
	if err := array.Set(teams); err != nil {
		return false, err
	}

	err = s.QueryRow(context.Background(), query, guildId, array).Scan(&valid)
	return
}

func (s *SupportTeamTable) Get(guildId uint64) (teams []SupportTeam, e error) {
	rows, err := s.Query(context.Background(), `SELECT "id", "name" from support_team WHERE "guild_id" = $1;`, guildId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		team := SupportTeam{
			GuildId: guildId,
		}

		if err := rows.Scan(&team.Id, &team.Name); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return
}

func (s *SupportTeamTable) GetWithMembers(guildId uint64) (teams map[SupportTeam][]uint64, e error) {
	query := `
SELECT support_team.id, support_team.name, support_team_members.user_id
FROM support_team
INNER JOIN support_team_members ON support_team.id = support_team_members.team_id
WHERE support_team.guild_id = $1;
`

	rows, err := s.Query(context.Background(), query, guildId)
	if err != nil {
		return nil, err
	}

	names := make(map[int]string)     // team_id -> name
	members := make(map[int][]uint64) // team_id -> [user_id]

	for rows.Next() {
		var teamId int
		var name string
		var userId uint64

		if err := rows.Scan(&teamId, &name, &userId); err != nil {
			return nil, err
		}

		if _, ok := names[teamId]; !ok {
			names[teamId] = name
		}

		if current, ok := members[teamId]; ok {
			members[teamId] = append(current, userId)
		} else {
			members[teamId] = make([]uint64, 0)
		}
	}

	teams = make(map[SupportTeam][]uint64)
	for id, name := range names {
		team := NewSupportTeam(id, guildId, name)
		teams[team] = members[id]
	}

	return
}

func (s *SupportTeamTable) Create(guildId uint64, name string) (id int, err error) {
	err = s.QueryRow(context.Background(), `INSERT INTO support_team("guild_id", "name") VALUES($1, $2) RETURNING "id";`, guildId, name).Scan(&id)
	return
}

func (s *SupportTeamTable) Delete(id int) (err error) {
	_, err = s.Exec(context.Background(), `DELETE FROM support_team WHERE "id"=$1;`, id)
	return
}
