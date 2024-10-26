.PHONY: build run clean

build:
	mkdir -p build
	go build -o ./build/server ./cmd/server
	go build -o ./build/client ./cmd/client

run: build
	./build/server

clean:
	rm -rf ./build