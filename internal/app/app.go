package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/raf-nr/xkcd-finder/internal/config"
	"github.com/raf-nr/xkcd-finder/internal/domain/comic"
	"github.com/raf-nr/xkcd-finder/internal/infra/xkcd"
	"github.com/raf-nr/xkcd-finder/internal/lib/language"
	"github.com/raf-nr/xkcd-finder/internal/lib/stemming"
	pgRepo "github.com/raf-nr/xkcd-finder/internal/repo/postgres"
	xkcdRepo "github.com/raf-nr/xkcd-finder/internal/repo/xkcd"
	"github.com/raf-nr/xkcd-finder/internal/usecase/search"
	"github.com/raf-nr/xkcd-finder/internal/usecase/upload"
)

type Application struct {
	UploadComicsUseCase  *upload.UploadComicsUseCase
	ReceiveComicsUseCase *search.ReceiveComicsUseCase
	db                   *sql.DB
}

func New(c *config.Config) (*Application, error) {
	db, err := sql.Open("postgres", c.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	languageDetector := language.NewLinguaGoDetector()
	stemmer := stemming.NewStemmer(languageDetector)
	comicGenerator := comic.NewComicGenerator(stemmer)

	comicRepo := pgRepo.NewComicRepository(db)

	httpClient := &http.Client{}
	comicService := xkcd.NewService(httpClient, c.XKCDService.BaseURL, c.XKCDService.ComicInfoPath)
	comicServiceRepo := xkcdRepo.NewComicServiceRepository(
		xkcdRepo.ComicConfiguration{Language: c.XKCDService.Language},
		comicService,
		comicGenerator,
	)

	receiveComicsUC := search.NewReceiveComicsUseCase(comicRepo, stemmer, c.XKCDService.Language)
	uploadComicsUC := upload.NewUploadComicsUseCase(comicServiceRepo, comicRepo, c.XKCDService.WorkersAmount)

	return &Application{
		UploadComicsUseCase:  uploadComicsUC,
		ReceiveComicsUseCase: receiveComicsUC,
		db:                   db,
	}, nil
}

func (a *Application) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}
