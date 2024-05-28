package comic

type Stemmer interface {
	Stem(string, string) []string
	StemWithLangDetection(string) ([]string, error)
}

type ComicGenerator struct {
	stemmer  Stemmer
	language string
}

func NewComicGenerator(stemmer Stemmer) *ComicGenerator {
	return &ComicGenerator{stemmer: stemmer}
}

func (cg *ComicGenerator) NewComic(imageURL string, info string) Comic {
	keywords := cg.stemmer.Stem(info, cg.language)
	return Comic{
		ImageURL: imageURL,
		Keywords: keywords,
	}
}

func (cg *ComicGenerator) NewComicWithLangDetection(imageURL string, info string) (Comic, error) {
	keywords, err := cg.stemmer.StemWithLangDetection(info)
	if err != nil {
		return Comic{}, err
	}

	return Comic{
		ImageURL: imageURL,
		Keywords: keywords,
	}, nil
}
