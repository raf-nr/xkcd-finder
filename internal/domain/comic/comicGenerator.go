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

func (cg *ComicGenerator) NewComic(imageURL string, info string, language string) (Comic, error) {
	if language == "" {
		return cg.newComicWithLangDetection(imageURL, info)
	}

	keywords := cg.stemmer.Stem(info, language)
	return Comic{
		ImageURL: imageURL,
		Keywords: keywords,
	}, nil
}

func (cg *ComicGenerator) newComicWithLangDetection(imageURL string, info string) (Comic, error) {
	keywords, err := cg.stemmer.StemWithLangDetection(info)
	if err != nil {
		return Comic{}, err
	}

	return Comic{
		ImageURL: imageURL,
		Keywords: keywords,
	}, nil
}
