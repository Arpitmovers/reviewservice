package main

import (
	"github.com/Arpitmovers/reviewservice/internal/app"
	"github.com/Arpitmovers/reviewservice/internal/config"
)

func main() {

	app := &app.App{}
	cfg := config.Load()
	app.Initialize(cfg)
	app.Run(":8080")

}
