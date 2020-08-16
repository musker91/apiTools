.PHONY: build clean tool lint help

all: build

build:
	@go build -v .
	@echo "build success"

pack:
	@go build -v .
	@bash ./build.sh
	@echo "packaged successfully"

pack-linux-amd64:
	@GOOS=linux GOARCH=amd64 go build -v .
	@bash ./build.sh
	@echo "packaged successfully"

tool:
	@go vet .
	@gofmt -w .

lint:
	@golint ./...

clean:
	@rm -rf apiTools
	@go clean -i .
	@rm -rf ./dist ./dist.zip

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
