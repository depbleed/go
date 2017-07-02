.PHONY: all build

all: build

build:
	go test -v ./go-depbleed
	go build -o bin/depbleed ./depbleed
