package main

import (
	"log"
	"os"
	"shadxel/internal/app"
	"shadxel/internal/config"
	"strings"
)

func main() {
	conf := parseConfig()
	shadxel, err := app.NewApp(conf)
	if err != nil {
		log.Fatalf("App crushed: %v", err)
	}
	shadxel.Run()
}

func parseConfig() config.Config {
	script := "demo"
	if len(os.Args) >= 2 {
		arg := os.Args[1]
		if strings.HasSuffix(arg, ".lua") {
			return config.Config{
				Script: arg,
			}
		}
		script = strings.TrimSuffix(os.Args[1], ".lua")
	}
	script = "demo/" + script + ".lua"

	return config.Config{
		Script: script,
	}
}
