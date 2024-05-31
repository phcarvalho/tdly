include .env
export

build:
	@CGO_ENABLED=1 go build -o bin/todo-app .

run: build
	@./bin/todo-app

migrate-up:
	goose up

migrate-down:
	goose down
