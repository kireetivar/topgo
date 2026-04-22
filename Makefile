.DEFAULT_GOAL := build

.PHONY : go vet build run
fmt: 
	go fmt ./...
vet:
	go vet ./...
build: fmt vet
	go build
run: fmt
	go run main.go