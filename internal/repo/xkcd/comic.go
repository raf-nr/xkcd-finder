package xkcd

import (
	"context"
	"fmt"

	"github.com/raf-nr/xkcd-finder/internal/domain/comic"
	"github.com/raf-nr/xkcd-finder/internal/infra/xkcd"
)

type ComicConfiguration struct {
	Language string
}

type ComicServiceRepository struct {
	service   *xkcd.Service
	generator *comic.ComicGenerator
	cfg       ComicConfiguration
}

func NewComicServiceRepository(cfg ComicConfiguration, service *xkcd.Service, generator *comic.ComicGenerator) *ComicServiceRepository {
	return &ComicServiceRepository{
		service:   service,
		generator: generator,
		cfg:       cfg,
	}
}

func (r *ComicServiceRepository) GetComic(ctx context.Context, id int) (comic.Comic, error) {
	comicRaw, err := r.service.GetComic(ctx, id)
	if err != nil {
		return comic.Comic{}, fmt.Errorf("fetch comic failed: %w", err)
	}

	comicFinal, err := r.generator.NewComic(comicRaw.Img, comicRaw.Description(), r.cfg.Language)
	if err != nil {
		return comic.Comic{}, fmt.Errorf("generate comic failed: %w", err)
	}

	return comicFinal, nil
}

func (r *ComicServiceRepository) GetComicsAmount(ctx context.Context) (int, error) {
	return r.service.GetComicsAmount(ctx)
}
