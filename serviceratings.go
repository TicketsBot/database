package database

import (
	"context"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceRatings struct {
	*pgxpool.Pool
}

func newServiceRatings(db *pgxpool.Pool) *ServiceRatings {
	return &ServiceRatings{
		db,
	}
}

func (ServiceRatings) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS service_ratings(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"rating" int2 NOT NULL,
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id"),
	PRIMARY KEY("guild_id", "ticket_id")
);`
}

func (r *ServiceRatings) Get(guildId uint64, ticketId int) (rating uint8, ok bool, e error) {
	query := `SELECT "rating" from service_ratings WHERE "guild_id" = $1 AND "ticket_id" = $2;`

	err := r.QueryRow(context.Background(), query, guildId, ticketId).Scan(&rating)
	if err == nil {
		return rating, true, nil
	} else if err == pgx.ErrNoRows {
		return 0, false, nil
 	} else {
 		return 0, false, err
	}
}

func (r *ServiceRatings) GetCount(guildId uint64) (count int, err error) {
	query := `SELECT COUNT(*) from service_ratings WHERE "guild_id" = $1;`
	err = r.QueryRow(context.Background(), query, guildId).Scan(&count)
	return
}

func (r *ServiceRatings) GetCountClaimedBy(guildId, userId uint64) (count int, err error) {
	query := `
SELECT COUNT(service_ratings.rating)
FROM service_ratings
INNER JOIN ticket_claims
ON service_ratings.guild_id = ticket_claims.guild_id AND service_ratings.ticket_id = ticket_claims.ticket_id
WHERE service_ratings.guild_id = $1 AND ticket_claims.user_id = $2;
`

	err = r.QueryRow(context.Background(), query, guildId).Scan(&count)
	return
}

// TODO: Materialized view?
func (r *ServiceRatings) GetAverage(guildId uint64) (average float32, err error) {
	query := `SELECT AVG(rating) from service_ratings WHERE "guild_id" = $1;`
	err = r.QueryRow(context.Background(), query, guildId).Scan(&average)
	return
}


// TODO: Materialized view?
func (r *ServiceRatings) GetAverageClaimedBy(guildId, userId uint64) (average float32, err error) {
	query := `
SELECT AVG(service_ratings.rating)
FROM service_ratings
INNER JOIN ticket_claims
ON service_ratings.guild_id = ticket_claims.guild_id AND service_ratings.ticket_id = ticket_claims.ticket_id
WHERE service_ratings.guild_id = $1 AND ticket_claims.user_id = $2;
`

	err = r.QueryRow(context.Background(), query, guildId, userId).Scan(&average)
	return
}

func (r *ServiceRatings) GetMulti(guildId uint64, ticketIds []uint) (map[int]uint8, error) {
	query := `SELECT "ticket_id", "rating" from service_ratings WHERE "guild_id" = $1 AND "ticket_id" = ANY($2);`

	idArray := &pgtype.Int4Array{}
	if err := idArray.Set(ticketIds); err != nil {
		return nil, err
	}

	ratings := make(map[int]uint8)

	rows, err := r.Query(context.Background(), query, guildId, idArray)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ticketId int
		var rating uint8

		if err := rows.Scan(&ticketId, &rating); err != nil {
			return nil, err
		}

		ratings[ticketId] = rating
	}

	return ratings, nil
}

// [lower,upper]
func (r *ServiceRatings) GetRange(guildId uint64, lowerId, upperId int) (map[int]uint8, error) {
	query := `SELECT "ticket_id", "rating" from service_ratings WHERE "guild_id" = $1 AND "ticket_id" >= $2 AND "ticket_id" <= 3;`

	ratings := make(map[int]uint8)

	rows, err := r.Query(context.Background(), query, guildId, lowerId, upperId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ticketId int
		var rating uint8

		if err := rows.Scan(&ticketId, &rating); err != nil {
			return nil, err
		}

		ratings[ticketId] = rating
	}

	return ratings, nil
}

func (r *ServiceRatings) Set(guildId uint64, ticketId int, rating uint8) (err error) {
	query := `
INSERT INTO service_ratings("guild_id", "ticket_id", "rating")
VALUES($1, $2, $3)
ON CONFLICT("guild_id", "ticket_id") DO UPDATE SET "rating" = 3;`

	_, err = r.Exec(context.Background(), query, guildId, ticketId, rating)
	return
}
