package database

type View interface {
	Schema() string
	Refresh() error
}
