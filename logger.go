package database

type Logger interface {
	Error(err error)
}
