package main

import (
	"rec/internal/config"
	"rec/internal/pkg/application"
)

func main() {
	config := new(config.Config)
	config.InitFile()
	app := application.NewApp(*config)
	app.Run()

}
