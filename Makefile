BIN := "./bin/fibonacci"

generate:
	protoc --proto_path=api/proto --go_out=. --go-grpc_out=. --grpc-gateway_out=. --validate_out="lang=go:." api/proto/Service.proto

test:
	go test -race -count 100 ./internal/...

teardown-integration-redis-tests:
	docker-compose -f docker-compose.redis.test.yml down

setup-integration-redis-tests: teardown-integration-redis-tests
	docker-compose -f docker-compose.redis.test.yml up -d

start-integration-redis-tests:
	GRPC_HOST="127.0.0.1" GRPC_PORT="9011" REDIS_HOST="localhost" REDIS_PORT=6379 REDIS_DB=1 go test -v -timeout=1s -count=1 ./tests/integration/redis -tags integration

start-bench-redis-tests:
	GRPC_HOST="127.0.0.1" GRPC_PORT="9011" REDIS_HOST="localhost" REDIS_PORT=6379 REDIS_DB=1 go test -timeout=1s -benchtime=100x -count=1 ./tests/integration/redis -tags bench

teardown-integration-memory-tests:
	docker-compose -f docker-compose.memory.test.yml down

setup-integration-memory-tests: teardown-integration-memory-tests
	docker-compose -f docker-compose.memory.test.yml up --build -d

start-integration-memory-tests:
	GRPC_HOST="127.0.0.1" GRPC_PORT="9011" go test -v -count=1 ./tests/integration/memory -tags integration

build:
	go build -v -o $(BIN) ./cmd

up:
	docker-compose up -d

down:
	docker-compose down

rebuild: down
	docker-compose up --build -d

lint:
	golangci-lint run ./...

.PHONY: up down
