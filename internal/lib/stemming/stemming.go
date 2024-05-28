package stemming

import (
	"fmt"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/raf-nr/xkcd-finder/internal/lib/iso6391"
)

type LanguageDetector interface {
	Detect(text string) (iso6391.ISO6391, error)
}

type Stemmer struct {
	langDetector LanguageDetector
}

func NewStemmer(detector LanguageDetector) *Stemmer {
	return &Stemmer{langDetector: detector}
}

func (s *Stemmer) StemWithLangDetection(text string) ([]string, error) {
	language, err := s.langDetector.Detect(text)
	if err != nil {
		return nil, fmt.Errorf("failed to detect language: %w", err)
	}
	return s.Stem(text, language), nil
}

func (s *Stemmer) Stem(text string, language iso6391.ISO6391) []string {
	words := strings.Fields(text)
	words = RemoveStopWords(words, language)
	return stemWords(words, language)
}

func stemWords(words []string, language iso6391.ISO6391) []string {
	result := make([]string, 0, len(words))

	for _, word := range words {
		stemmed, err := snowball.Stem(word, language.Name(), true)
		if err != nil {
			result = append(result, word) // fallback: use original word if stemming fails
			continue
		}
		result = append(result, stemmed)
	}

	return removeDuplicates(result)
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
