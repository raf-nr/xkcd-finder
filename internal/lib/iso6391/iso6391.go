package iso6391

import (
	"errors"
	"strings"

	"github.com/pemistahl/lingua-go"
)

type ISO6391 struct {
	code string
}

func NewISO6391(code string) (ISO6391, error) {
	// TODO: Validate the ISO 639-1 code
	return ISO6391{code: code}, nil
}

func DetectISO6391FromText(text string) (ISO6391, error) {
	detector := lingua.NewLanguageDetectorBuilder().FromAllLanguages().Build()

	lang, ok := detector.DetectLanguageOf(text)
	if !ok {
		return ISO6391{}, errors.New("language detection failed")
	}

	isoCode := lang.IsoCode639_1().String()
	if isoCode == "" {
		return ISO6391{}, errors.New("no ISO 639-1 code for detected language")
	}

	return NewISO6391(isoCode)
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
