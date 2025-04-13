package main

import (
	"fmt"
	"log"
	"os"
	"shadxel/internal/app"
	"shadxel/internal/config"
	"strings"
)

func main() {
	conf, err := parseConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	shadxel, err := app.NewApp(*conf)
	if err != nil {
		log.Fatalf("App crushed: %v", err)
	}
	shadxel.Run()
}

func parseConfig() (*config.Config, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("usage: %s <script.lua>", os.Args[0])
	}
	script := strings.TrimSuffix(os.Args[1], ".lua")
	script += ".lua"
	script = "demo/" + script

	return &config.Config{
		Script: script,
	}, nil
}
