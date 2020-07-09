package database

type MessageId int

const (
	NoPermission MessageId = iota
)

var Messages = map[string]MessageId{
	"no_permission": NoPermission,
}