build:
	go build -o bin/shadxel cmd/main.go
run: build
	./bin/shadxel

devrun:
	go run cmd/main.go menger -size=5 -cpu=8

win:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
	go build -o build/shadxel.exe ./cmd/main.go
