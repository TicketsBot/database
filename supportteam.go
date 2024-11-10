package database

import (
	"context"
	"errors"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SupportTeamTable struct {
	*pgxpool.Pool
}

type SupportTeam struct {
	Id         int     `json:"id"`
	GuildId    uint64  `json:"guild_id"`
	Name       string  `json:"name"`
	OnCallRole *uint64 `json:"on_call_role_id"`
}

func NewSupportTeam(id int, guildId uint64, name string, onCallRole *uint64) SupportTeam {
	return SupportTeam{
		Id:         id,
		GuildId:    guildId,
		Name:       name,
		OnCallRole: onCallRole,
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
	"on_call_role_id" int8 DEFAULT NULL UNIQUE,
	UNIQUE("guild_id", "name"),
	PRIMARY KEY("id")
);`
}

func (s *SupportTeamTable) Exists(ctx context.Context, teamId int, guildId uint64) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM support_team WHERE "id" = $1 and "guild_id" = $2);`
	err = s.QueryRow(ctx, query, teamId, guildId).Scan(&exists)
	return
}

func (s *SupportTeamTable) AllTeamsMatchGuild(ctx context.Context, guildId uint64, teams []int) (valid bool, err error) {
	query := `SELECT NOT EXISTS(SELECT 1 FROM support_team WHERE "id" = ANY($1) and "guild_id" != $2);`

	array := &pgtype.Int4Array{}
	if err := array.Set(teams); err != nil {
		return false, err
	}

	err = s.QueryRow(ctx, query, array, guildId).Scan(&valid)
	return
}

func (s *SupportTeamTable) AllTeamsExistForGuild(ctx context.Context, guildId uint64, teams []int) (valid bool, err error) {
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

	err = s.QueryRow(ctx, query, guildId, array).Scan(&valid)
	return
}

func (s *SupportTeamTable) Get(ctx context.Context, guildId uint64) (teams []SupportTeam, e error) {
	rows, err := s.Query(ctx, `SELECT "id", "name", "on_call_role_id" from support_team WHERE "guild_id" = $1;`, guildId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		team := SupportTeam{
			GuildId: guildId,
		}

		if err := rows.Scan(&team.Id, &team.Name, &team.OnCallRole); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return
}

func (s *SupportTeamTable) GetById(ctx context.Context, guildId uint64, id int) (SupportTeam, bool, error) {
	query := `SELECT "name", "on_call_role_id" from support_team WHERE "guild_id" = $1 AND "id" = $2;`

	team := SupportTeam{
		Id:      id,
		GuildId: guildId,
	}

	if err := s.QueryRow(ctx, query, guildId, id).Scan(&team.Name, &team.OnCallRole); err != nil {
		if err == pgx.ErrNoRows {
			return SupportTeam{}, false, nil
		} else {
			return SupportTeam{}, false, err
		}
	}

	return team, true, nil
}

func (s *SupportTeamTable) GetByName(ctx context.Context, guildId uint64, name string) (SupportTeam, bool, error) {
	query := `SELECT "id", "name", "on_call_role_id" from support_team WHERE "guild_id" = $1 AND "name" = $2;`

	team := SupportTeam{
		GuildId: guildId,
	}

	if err := s.QueryRow(ctx, query, guildId, name).Scan(&team.Id, &team.Name, &team.OnCallRole); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SupportTeam{}, false, nil
		} else {
			return SupportTeam{}, false, err
		}
	}

	return team, true, nil
}

func (s *SupportTeamTable) GetMulti(ctx context.Context, guildId uint64, teamIds []int) (map[int]SupportTeam, error) {
	query := `
SELECT "id", "name", "on_call_role_id"
FROM support_team
WHERE "guild_id" = $1 AND "id" = ANY($2);`

	array := &pgtype.Int4Array{}
	if err := array.Set(teamIds); err != nil {
		return nil, err
	}

	rows, err := s.Query(ctx, query, guildId, array)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	teams := make(map[int]SupportTeam)
	for rows.Next() {
		team := SupportTeam{
			GuildId: guildId,
		}

		if err := rows.Scan(&team.Id, &team.Name, &team.OnCallRole); err != nil {
			return nil, err
		}

		teams[team.Id] = team
	}

	return teams, nil
}

func (s *SupportTeamTable) GetWithMembers(ctx context.Context, guildId uint64) (teams map[SupportTeam][]uint64, e error) {
	query := `
SELECT support_team.id, support_team.name, support_team.on_call_role_id, support_team_members.user_id
FROM support_team
INNER JOIN support_team_members ON support_team.id = support_team_members.team_id
WHERE support_team.guild_id = $1;
`

	rows, err := s.Query(ctx, query, guildId)
	if err != nil {
		return nil, err
	}

	names := make(map[int]string)        // team_id -> name
	onCallRoles := make(map[int]*uint64) // team_id -> on_call_role_id
	members := make(map[int][]uint64)    // team_id -> [user_id]

	for rows.Next() {
		var teamId int
		var name string
		var onCallRoleId *uint64
		var userId uint64

		if err := rows.Scan(&teamId, &name, &onCallRoleId, &userId); err != nil {
			return nil, err
		}

		if _, ok := names[teamId]; !ok {
			names[teamId] = name
		}

		if _, ok := onCallRoles[teamId]; !ok {
			onCallRoles[teamId] = onCallRoleId
		}

		if current, ok := members[teamId]; ok {
			members[teamId] = append(current, userId)
		} else {
			members[teamId] = make([]uint64, 0)
		}
	}

	teams = make(map[SupportTeam][]uint64)
	for id, name := range names {
		team := NewSupportTeam(id, guildId, name, onCallRoles[id])
		teams[team] = members[id]
	}

	return
}

func (s *SupportTeamTable) Create(ctx context.Context, guildId uint64, name string) (id int, err error) {
	err = s.QueryRow(ctx, `INSERT INTO support_team("guild_id", "name") VALUES($1, $2) RETURNING "id";`, guildId, name).Scan(&id)
	return
}

func (s *SupportTeamTable) SetOnCallRole(ctx context.Context, teamId int, roleId *uint64) (err error) {
	_, err = s.Exec(ctx, `UPDATE support_team SET "on_call_role_id" = $2 WHERE "id" = $1;`, teamId, roleId)
	return
}

func (s *SupportTeamTable) Delete(ctx context.Context, id int) (err error) {
	_, err = s.Exec(ctx, `DELETE FROM support_team WHERE "id"=$1;`, id)
	return
}
