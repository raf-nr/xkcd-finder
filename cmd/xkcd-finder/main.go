package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/raf-nr/xkcd-finder/internal/app"
	"github.com/raf-nr/xkcd-finder/internal/config"
)

func main() {
	config, err := config.New(config.DefaultPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app, err := app.New(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer app.Close()
}
