.PHONY: build clean tool lint help

all: build

build:
	@go build -v .
	@echo "build success"

tool:
	@go vet .
	@gofmt -w .

lint:
	@golint ./...

clean:
	@rm -rf apiTools
	@go clean -i .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"