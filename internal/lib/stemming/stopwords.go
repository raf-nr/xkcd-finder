package stemming

import (
	"github.com/raf-nr/xkcd-finder/internal/lib/iso6391"
	stopwordsiso "github.com/toadharvard/stopwords-iso"
)

func RemoveStopWords(text []string, language iso6391.ISO6391) []string {
	stopwordsMapping, _ := stopwordsiso.NewStopwordsMapping()

	var result []string
	for _, word := range text {
		cleanWord := stopwordsMapping.ClearStringByLang(word, language.Code())
		if cleanWord != "" {
			result = append(result, cleanWord)
		}
	}
	return result
}
