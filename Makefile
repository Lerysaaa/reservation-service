.PHONY: test run docker-up docker-down

test:
	go test ./...

run:
	go run ./cmd/api

docker-up:
	docker compose up --build

docker-down:
	docker compose down
