package search

import (
	"context"

	"github.com/raf-nr/xkcd-finder/internal/domain/comic"
)

type ComicRepository interface {
	FindTopComicsByKeywords(ctx context.Context, keywords []string, limit uint) ([]comic.Comic, error)
}

type Stemmer interface {
	Stem(text string, lang string) []string
	StemWithLangDetection(text string) ([]string, error)
}

type ReceiveComicsUseCase struct {
	comicRepo ComicRepository
	stemmer   Stemmer
	language  string
}

func NewReceiveComicsUseCase(comicRepo ComicRepository, stemmer Stemmer, language string) *ReceiveComicsUseCase {
	return &ReceiveComicsUseCase{
		comicRepo: comicRepo,
		stemmer:   stemmer,
		language:  language,
	}
}

func (uc *ReceiveComicsUseCase) Execute(ctx context.Context, query string, limit uint) ([]comic.Comic, error) {
	stemmedKeywords := uc.stemmer.Stem(query, uc.language)

	comics, err := uc.comicRepo.FindTopComicsByKeywords(ctx, stemmedKeywords, limit)
	if err != nil {
		return nil, err
	}

	return comics, nil
}
