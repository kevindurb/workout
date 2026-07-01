export GOOSE_DRIVER := "sqlite"
export GOOSE_DBSTRING := "./database.db"
export GOOSE_MIGRATION_DIR := "./internal/database/migrations"

default:
  @just --list

[parallel]
dev: air sql-watch

run:
  go run cmd/planner/main.go

test:
  go test -v ./...

air:
  go tool air

test-watch:
  watchexec -e go,tmpl,sql just test

e2e-watch:
  watchexec -e go,sql,hurl just e2e

sql-watch:
  watchexec -e sql just sqlc generate

sqlfluff-fix:
  sqlfluff fix ./internal/database/**/*.sql

pre-commit-install:
  pre-commit install

goose *ARGS:
  go tool goose {{ARGS}}

sqlc *ARGS:
  go tool sqlc {{ARGS}}
