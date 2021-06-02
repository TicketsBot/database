package database

type Language string

const (
	English       Language = "en"
	French        Language = "fr"
	Spanish       Language = "es"
	German        Language = "de"
	Dutch         Language = "nl"
	Polish        Language = "pl"
	Norwegian     Language = "no"
	Turkish       Language = "tr"
	Swedish       Language = "sv"
	Arabic        Language = "ar"
	Hungarian     Language = "hu"
	Russian       Language = "ru"
	PortugueseBR  Language = "br"
	Chinese       Language = "cn"
	ChineseTaiwan Language = "tw"
)

var Flags = map[Language]string{
	English:       "🇬🇧",
	French:        "🇫🇷",
	Spanish:       "🇪🇸",
	German:        "🇩🇪",
	Dutch:         "🇳🇱",
	Polish:        "🇵🇱",
	Norwegian:     "🇳🇴",
	Turkish:       "🇹🇷",
	Swedish:       "🇸🇪",
	Arabic:        "🇸🇦",
	Hungarian:     "🇭🇺",
	Russian:       "🇷🇺",
	PortugueseBR:  "🇧🇷",
	Chinese:       "🇨🇳",
	ChineseTaiwan: "🇹🇼",
}

// https://discord.com/developers/docs/dispatch/field-values
var Locales = map[string]Language{
	"en-US": English,
	"en-GB": English,
	"zh-CN": Chinese,
	"zh-TW": ChineseTaiwan,
	"cs":    English, // Czech
	"da":    English, // Danish
	"nl":    Dutch,
	"fr":    French,
	"de":    German,
	"hu":    Hungarian,
	"it":    English, // Italian
	"ja":    English, // Japanese
	"ko":    English, // Korean
	"no":    Norwegian,
	"pl":    Polish,
	"pt-BR": PortugueseBR,
	"ru":    Russian,
	"es-ES": Spanish,
	"sv-SE": Swedish,
	"tr":    Turkish,
	"bg":    English, // Bulgarian
	"uk":    English, // Ukrainian
	"fi":    English, // Finnish
	"hr":    English, // Croatian
	"ro":    English, // Romanian
	"lt":    English, // Lithuanian
}
