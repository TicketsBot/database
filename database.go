package database

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	*pgxpool.Pool
	Logger Logger
}

func NewDatabase(pool *pgxpool.Pool, logger Logger) *Database {
	return &Database{
		Pool: pool,
		Logger: logger,
	}
}
