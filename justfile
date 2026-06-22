db_file := "database.db"

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
  hurl -v test/e2e

sqlfluff-fix:
  sqlfluff fix ./internal/database/**/*.sql

pre-commit-install:
  pre-commit install
