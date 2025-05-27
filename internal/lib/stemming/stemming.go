package stemming

import (
	"fmt"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/raf-nr/xkcd-finder/internal/lib/iso6391"
)

type Stemmer struct {
	text     string
	language iso6391.ISO6391
}

func NewStemmer(text string) (Stemmer, error) {
	lang, err := iso6391.DetectISO6391FromText(text)
	if err != nil {
		return Stemmer{}, fmt.Errorf("failed to detect language: %s", err)
	}
	return Stemmer{
		text:     text,
		language: lang,
	}, nil
}

func (s *Stemmer) Stem() []string {
	words := strings.Fields(s.text)
	words = RemoveStopWords(words, s.language)
	words = stemText(words, s.language)
	words = removeDuplicates(words)
	return words
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
