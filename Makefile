.PHONY:

build-graphql-model:
	go run github.com/99designs/gqlgen generate
# go run -mod=mod github.com/99designs/gqlgen generate --verbose
install:
	go get ./...

prepare: install migrate-up

start:
	go run cmd/main.go

test:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./tests/...

cover:
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

build:
	go build -o main ./cmd/main.go

docker-build:
	docker build -t admin-franchise:v1 -f docker/Dockerfile .

docker-compose-up:
	docker-compose -f ./docker/docker-compose.yml up --build

create-migration:
	migrate create -ext sql -dir setup/migrations -seq $(name)