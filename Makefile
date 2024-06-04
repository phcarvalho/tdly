include .env
export

build:
	@CGO_ENABLED=1 go build -o bin/tdly .

run: build
	@./bin/tdly

migrate-up:
	goose up

migrate-down:
	goose down
