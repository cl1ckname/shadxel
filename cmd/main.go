package main

import (
	"log"
	"shadxel/internal/app"
	"shadxel/internal/config"
)

func main() {
	conf := config.ParseConfig()
	shadxel, err := app.NewApp(conf)
	if err != nil {
		log.Fatalf("App crushed: %v", err)
	}
	shadxel.Run()
}
