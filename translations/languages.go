package database

type Language string

const (
	English Language = "en"
)

var Languages = map[Language]string{
	English: "🇬🇧",
}
