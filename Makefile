.PHONY: build run kill clean

build:
	mkdir -p build/master/logs
	mkdir -p build/replica1/logs
	mkdir -p build/clientdir/logs

	go build -o ./build/server ./cmd/server
	go build -o ./build/client ./cmd/client

	cp ./build/server ./build/master/
	cp ./build/server ./build/replica1/
	cp ./build/client ./build/clientdir/

kill:
	- if [ -f build/master/server.pid ]; then kill -9 $$(cat build/master/server.pid); fi
	- if [ -f build/replica1/server.pid ]; then kill -9 $$(cat build/replica1/server.pid); fi

run: build kill
	./build/master/server --config ./build/master/go-redis.conf & echo $$! > build/master/server.pid
	./build/replica1/server --config ./build/replica1/go-redis.conf & echo $$! > build/replica1/server.pid

clean:
	rm -rf ./build