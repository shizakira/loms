include .env
export $(shell sed 's/=.*//' .env)

up:
	docker compose up --build --force-recreate

down:
	docker compose down

grpc_gen:
	buf dep update
	mkdir -p ./gen/grpc
	buf generate

grpc-dep-update:
	buf dep update
	buf build

grpc-dep-install:
	grpc-dep

mockery-install:
	go install github.com/vektra/mockery/v3@v3.2.5

.PHONY: test
test:
	go test -v -cover ./...

integration-test:
	go test -count=1 -v -tags=integration ./test/integration

sqlc-install:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

sqlc-generate:
	sqlc generate

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all