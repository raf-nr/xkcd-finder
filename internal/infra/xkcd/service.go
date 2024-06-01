package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type configuration struct {
	baseURL       string
	comicInfoPath string
}

type Service struct {
	client Client
	configuration
}

func NewService(client Client, baseURL string, comicInfoPath string) *Service {
	return &Service{
		client: client,
		configuration: configuration{
			baseURL:       baseURL,
			comicInfoPath: comicInfoPath,
		},
	}
}

func (c *Service) GetComic(ctx context.Context, id int) (Comic, error) {
	url := fmt.Sprintf("%s/%d/%s", c.baseURL, id, c.comicInfoPath)
	return c.fetchComic(ctx, url)
}

func (c *Service) GetComicsAmount(ctx context.Context) (int, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, c.comicInfoPath)
	comic, err := c.fetchComic(ctx, url)
	if err != nil {
		return 0, err
	}
	return comic.Num, nil
}

func (c *Service) fetchComic(ctx context.Context, url string) (Comic, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Comic{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Comic{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Comic{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return Comic{}, fmt.Errorf("failed to decode comic JSON: %w", err)
	}

	return comic, nil
}
