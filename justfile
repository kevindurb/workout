export GOOSE_DRIVER := "sqlite"
export GOOSE_DBSTRING := "./database.db"
export GOOSE_MIGRATION_DIR := "./internal/database/migrations"

default:
  @just --list

run:
  go run cmd/planner/main.go

test:
  go test -v ./...

dev:
  air

test-watch:
  watchexec -e go,tmpl,sql just test

e2e:
  hurl test/e2e

e2e-watch:
  watchexec -e go,sql,hurl just e2e

sqlfluff-fix:
  sqlfluff fix ./internal/database/**/*.sql

pre-commit-install:
  pre-commit install

goose *ARGS:
  goose {{ARGS}}

sqlc *ARGS:
  sqlc {{ARGS}}
