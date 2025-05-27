package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raf-nr/xkcd-finder/internal/domain"
	"github.com/raf-nr/xkcd-finder/internal/lib/stemming"
)

type Comic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

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

func (c *Service) GetComic(id int) (domain.Comic, error) {
	url := fmt.Sprintf("%s/%d/%s", c.baseURL, id, c.comicInfoPath)
	comic, err := c.fetchComic(url)
	if err != nil {
		return domain.Comic{}, err
	}

	stemmer, err := stemming.NewStemmer(comic.Transcript)
	if err != nil {
		return domain.Comic{}, fmt.Errorf("failed to create stemmer: %w", err)
	}

	return domain.Comic{
		URL:      comic.Img,
		Keywords: stemmer.Stem(),
	}, nil
}

func (c *Service) GetLastComicID() (int, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, c.comicInfoPath)
	comic, err := c.fetchComic(url)
	if err != nil {
		return 0, err
	}
	return comic.Num, nil
}

func (c *Service) fetchComic(url string) (Comic, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
