package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var DefaultPath = filepath.Join(findModuleRootOrFail(), "config", "config.yaml")

type Postgres struct {
	DSN string `yaml:"dsn"`
}

type XKCDService struct {
	BaseURL       string `yaml:"base-url"`
	ComicInfoPath string `yaml:"comic-info-path"`
	WorkersAmount int    `yaml:"workers-amount"`
	Language      string `yaml:"language"`
}

type HTTPServer struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type Config struct {
	Postgres    Postgres    `yaml:"postgres"`
	XKCDService XKCDService `yaml:"xkcd-service"`
	HTTPServer  HTTPServer  `yaml:"http-server"`
}

func New(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func findModuleRootOrFail() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("go.mod not found")
		}
		dir = parent
	}
}
