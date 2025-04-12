package main

import (
	"log"
	"shadxel/internal/app"
)

func main() {
	shadxel, err := app.NewApp()
	if err != nil {
		log.Fatalf("App crushed: %v", err)
	}
	shadxel.Run()
}
