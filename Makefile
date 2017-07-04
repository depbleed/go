.PHONY: all build

all: build

build:
	go test -v ./go-depbleed -coverprofile=coverage.txt
	go build -o bin/depbleed ./depbleed