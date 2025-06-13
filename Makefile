.PHONY: run build test docker-up docker-down clean

run:
	go run cmd/server/main.go

build:
	mkdir -p build
	go build -o build/server cmd/server/main.go

test:
	go test ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

clean:
	rm -rf build
