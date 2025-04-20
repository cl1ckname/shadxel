build:
	go build -o bin/shadxel cmd/main.go
run: build
	./bin/shadxel

devrun:
	go run cmd/main.go lua/template.lua
