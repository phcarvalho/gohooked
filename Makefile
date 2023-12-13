dev:
	@go run .

build:
	@go build -o ./bin/gohooked

start:
	@./bin/gohooked
