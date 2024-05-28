package iso6391

import (
	"strings"
)

type ISO6391 struct {
	code string
}

func NewISO6391(code string) (ISO6391, error) {
	// TODO: Validate the ISO 639-1 code
	return ISO6391{code: code}, nil
}

func (i *ISO6391) Code() string {
	return i.code
}

func (i *ISO6391) Name() string {
	switch strings.ToLower(i.code) {
	case "en":
		return "english"
	case "es":
		return "spanish"
	case "fr":
		return "french"
	case "de":
		return "german"
	case "zh":
		return "chinese"
	case "ja":
		return "japanese"
	case "ru":
		return "Russian"
	case "ar":
		return "arabic"
	default:
		return "english"
	}
}
