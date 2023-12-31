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
	docker rmi -f bond-trading
	docker build -t bond-trading:latest .
	just clean-docker-cache

clean-docker-cache:
    docker builder prune -af
