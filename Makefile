
test:
	@go test ./... -race

build:
	@docker build .
