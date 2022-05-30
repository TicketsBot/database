package database

type View interface {
	Refresh() error
}
