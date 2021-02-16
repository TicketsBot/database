package database

import (
	"github.com/jackc/pgtype"
	"time"
)

func toInterval(duration time.Duration) (interval pgtype.Interval, err error) {
	err = interval.Set(duration)
	return
}
