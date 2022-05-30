package main

import (
	"context"
	"github.com/TicketsBot/database"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	logrus.Info("Connecting to database...")
	pool := must(pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URI")))
	db := database.NewDatabase(pool)
	logrus.Info("Connected!")

	if os.Getenv("DAEMON") == "true" {
		for {
			doRefresh(db)
			time.Sleep(6 * time.Hour)
		}
	} else {
		doRefresh(db)
	}
}

func doRefresh(db *database.Database) {
	logrus.Info("Starting refresh...")

	for _, view := range db.Views() {
		if err := view.Refresh(); err != nil {
			logrus.Errorf("Error refreshing view: %s", err.Error())
		}
	}

	logrus.Info("Refresh complete")
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
