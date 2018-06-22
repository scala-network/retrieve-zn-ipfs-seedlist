#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := retrieve-zn-ipfs-seedlist

default: build

build_linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -o ./bin/${APP_NAME}-linux-amd64 ./src/main.go

build_windows:
	GOOS=windows \
	GOARCH=amd64 \
	go build -o ./bin/${APP_NAME}-win-amd64.exe ./src/main.go

build_macos:
	GOOS=darwin \
	GOARCH=amd64 \
	go build -o ./bin/${APP_NAME}-macos-amd64 ./src/main.go

build: build_linux \
	build_windows \
	build_macos

run: build_linux
	LOG_FORMAT=Text \
	LOG_LEVEL=Debug \
	./bin/${APP_NAME}-linux-amd64

run_race:
	LOG_OUTPUT=Text \
	LOG_LEVEL=Debug \
	go run -race ./src/main.go

clean:
	rm ./bin/*
