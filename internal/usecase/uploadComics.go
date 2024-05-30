package usecase

import (
	"context"
	"sync"

	"github.com/raf-nr/xkcd-finder/internal/domain/comic"
)

type ComicServiceRepository interface {
	GetComic(ctx context.Context, id int) (comic.Comic, error)
	GetComicsAmount(ctx context.Context) (int, error)
}

type ComicRepository interface {
	Save(ctx context.Context, comic comic.Comic) error
	GetComics(ctx context.Context) ([]comic.Comic, error)
	GetComicsIDs(ctx context.Context) ([]int, error)
}

type UploadComicsUseCase struct {
	comicServiceRepo ComicServiceRepository
	comicRepo        ComicRepository
	workersAmount    int
}

func NewUploadComicsUseCase(
	comicServiceRepo ComicServiceRepository,
	comicRepo ComicRepository,
	workersAmount int,
) *UploadComicsUseCase {
	return &UploadComicsUseCase{
		comicServiceRepo: comicServiceRepo,
		comicRepo:        comicRepo,
		workersAmount:    workersAmount,
	}
}

func (uc *UploadComicsUseCase) Execute(ctx context.Context) error {
	comicsAmount, err := uc.comicServiceRepo.GetComicsAmount(ctx)
	if err != nil {
		return err
	}

	existingComics, err := uc.comicRepo.GetComicsIDs(ctx)
	if err != nil {
		return err
	}

	missingIDs := findMissingIDs(existingComics, comicsAmount)

	idChan := make(chan int)
	wg := sync.WaitGroup{}

	for i := 1; i <= uc.workersAmount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idChan {
				comicData, err := uc.comicServiceRepo.GetComic(ctx, id)
				if err != nil {
					continue // TODO: log error
				}
				if err := uc.comicRepo.Save(ctx, comicData); err != nil {
					continue // TODO: log error
				}
			}
		}()
	}

	go func() {
		defer close(idChan)
		for _, id := range missingIDs {
			idChan <- id
		}
	}()

	wg.Wait()
	return nil
}

func findMissingIDs(existing []int, maxID int) []int {
	existingSet := make(map[int]struct{}, len(existing))
	for _, id := range existing {
		existingSet[id] = struct{}{}
	}

	var missing []int
	for i := 1; i <= maxID; i++ {
		if _, ok := existingSet[i]; !ok {
			missing = append(missing, i)
		}
	}
	return missing
}
