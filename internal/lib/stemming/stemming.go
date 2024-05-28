package stemming

import (
	"fmt"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/raf-nr/xkcd-finder/internal/lib/iso6391"
)

type Stemmer struct {
	language *iso6391.ISO6391
}

func (s *Stemmer) StemWithLangDetection(text string) ([]string, error) {
	lang, err := iso6391.DetectISO6391FromText(text)
	if err != nil {
		return nil, fmt.Errorf("failed to detect language: %s", err)
	}
	s.language = &lang

	return s.Stem(text), nil
}

func (s *Stemmer) Stem(text string) []string {
	words := RemoveStopWords(strings.Fields(text), *s.language)
	words = stemText(words, *s.language)
	return removeDuplicates(words)
}

func stemText(words []string, language iso6391.ISO6391) []string {
	stemmedWords := make([]string, 0, len(words))

	for _, word := range words {
		if stemmed, err := snowball.Stem(word, language.Name(), true); err != nil {
			stemmedWords = append(stemmedWords, word)
		} else {
			stemmedWords = append(stemmedWords, stemmed)
		}
	}

	return removeDuplicates(stemmedWords)
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
