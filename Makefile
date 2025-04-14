build:
	go build -o bin/shadxel cmd/main.go
run: build
	./bin/shadxel
