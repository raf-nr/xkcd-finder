package detector

import (
	"errors"

	"github.com/pemistahl/lingua-go"
	"github.com/raf-nr/xkcd-finder/internal/lib/iso6391"
)

type LinguaGoDetector struct {
	detector lingua.LanguageDetector
}

func NewLinguaGoDetector() *LinguaGoDetector {
	detector := lingua.NewLanguageDetectorBuilder().FromAllLanguages().Build()
	return &LinguaGoDetector{detector: detector}
}

func (l *LinguaGoDetector) Detect(text string) (iso6391.ISO6391, error) {
	lang, ok := l.detector.DetectLanguageOf(text)
	if !ok {
		return iso6391.ISO6391{}, errors.New("language detection failed")
	}

	isoCode := lang.IsoCode639_1().String()
	if isoCode == "" {
		return iso6391.ISO6391{}, errors.New("no ISO 639-1 code for detected language")
	}

	return iso6391.NewISO6391(isoCode)
}
