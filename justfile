set windows-shell := ["powershell", "-NoLogo", "-Command"]

default:
    just --list

test:
	go test -race -short $$(go list ./...)

lint:
	golangci-lint run ./...

build:
	CGO_ENABLED=0 go build -o ./cmd/trading-server ./cmd

clean-go-cache:
	go clean -cache -modcache -i -r

build-docker:
	docker build -t api-service:latest -f ./docker/api_service.dockerfile ./docker/
	docker build -t order-service:latest -f ./docker/order_service.dockerfile ./docker/
	docker image prune
	just clean-docker-cache

clean-docker-cache:
    docker builder prune -af

compose-up:
	docker compose -f ./docker/all.compose.yml -p bond-trading up -d

memphis-up:
	docker compose -f ./docker/memphis.compose.yml -p memphis up -d

clickhouse-up:
	docker compose -f ./docker/clickhouse.compose.yml -p clickhouse up -d

tsdb-up:
	docker compose -f ./docker/timescaledb.compose.yml -p timescaledb up -d
