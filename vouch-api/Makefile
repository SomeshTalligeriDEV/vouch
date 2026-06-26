.PHONY: dev build test lint docker-build docker-up docker-down clean

# ── Variables ──────────────────────────────────────────────────────────────────
BINARY_API    := bin/api
BINARY_WORKER := bin/worker
MODULE        := github.com/SomeshTalligeriDEV/vouch-api

# ── Development ────────────────────────────────────────────────────────────────

## dev: run the API server with hot-reload via air (install: go install github.com/air-verse/air@latest)
dev:
	air -c .air.toml

## build: compile both API and worker binaries into ./bin
build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_API)    ./cmd/api
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_WORKER) ./cmd/worker

## test: run all unit tests with race detector
test:
	go test -race -count=1 ./...

## lint: run golangci-lint (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...

# ── Docker ─────────────────────────────────────────────────────────────────────

## docker-build: build all images defined in docker/docker-compose.yml
docker-build:
	docker compose -f docker/docker-compose.yml build

## docker-up: start the full local stack (mongo + redis + api + worker)
docker-up:
	docker compose -f docker/docker-compose.yml up -d

## docker-down: stop and remove containers
docker-down:
	docker compose -f docker/docker-compose.yml down

# ── Housekeeping ───────────────────────────────────────────────────────────────

## clean: remove built binaries
clean:
	rm -rf bin/

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
