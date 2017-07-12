.PHONY: all build test

all: build test

build:
	go build -o bin/depbleed ./depbleed

test:
	go test ./go-depbleed -coverprofile=coverage.txt
	go tool cover -func=coverage.txt
