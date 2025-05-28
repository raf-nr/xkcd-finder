package comic

type Stemmer interface {
	Stem(string, string) []string
	StemWithLangDetection(string) ([]string, error)
}

type ComicGenerator struct {
	stemmer Stemmer
}

func NewComicGenerator(stemmer Stemmer) *ComicGenerator {
	return &ComicGenerator{stemmer: stemmer}
}

func (cg *ComicGenerator) NewComic(id int, imageURL string, info string, language string) (Comic, error) {
	if language == "" {
		return cg.newComicWithLangDetection(id, imageURL, info)
	}

	keywords := cg.stemmer.Stem(info, language)
	return Comic{
		ID:       id,
		ImageURL: imageURL,
		Keywords: keywords,
	}, nil
}

func (cg *ComicGenerator) newComicWithLangDetection(id int, imageURL string, info string) (Comic, error) {
	keywords, err := cg.stemmer.StemWithLangDetection(info)
	if err != nil {
		return Comic{}, err
	}

	return Comic{
		ID:       id,
		ImageURL: imageURL,
		Keywords: keywords,
	}, nil
}
