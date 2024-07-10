package database

import "context"

type View interface {
	Refresh(ctx context.Context) error
}
