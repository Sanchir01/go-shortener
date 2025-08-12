SILENT:
PHONY:
MIGRATION_NAME ?= new_migration

DB_CONN_DEV = "host=localhost user=postgres password=postgres port=5443 dbname=shortener sslmode=disable"
DB_CONN_PROD = "host=103.74.92.37 user=postgres password=VpfthnhDDrvyH2Yi9J0b6m38udda2e918Tlm6JORwhjeS8R4ko6sM9HJmGdo2qEW port=5444 dbname=postgres"
FOLDER_PG= migrations/pg

compose:
	docker-compose up -d


build:
	go build -o ./.bin/main ./cmd/main/main.go
run: build
	./.bin/main
swag:
	swag init -g cmd/main/main.go
migrations-click-up:
	goose -dir migrations clickhouse "tcp://localhost:9000?username=default&password=clickhouse" up

lint:
	golangci-lint run

tests:
	go test ./test/...

migrations-up:
	goose -dir $(FOLDER_PG) postgres $(DB_CONN_DEV)   up

migrations-up-prod:
	goose -dir $(FOLDER_PG) postgres $(DB_CONN_PROD)   up

migrations-down:
	goose -dir $(FOLDER_PG) postgres $(DB_CONN_DEV)   down


migrations-status:
	goose -dir $(FOLDER_PG) postgres $(DB_CONN_DEV)  status

migrations-new:
	goose -dir $(FOLDER_PG) create $(MIGRATION_NAME) sql

gen:
	go generate ./...

docker:
	docker compose up -d

compose-prod:
	docker compose -f docker-compose.prod.yaml up --build -d
