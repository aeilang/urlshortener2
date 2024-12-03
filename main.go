package main

import "github.com/aeilang/urlshortener/app"

func main() {
	app, err := app.NewApplication("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	app.Run()
}
