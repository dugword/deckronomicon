# Makefile for Tron Eggs Simulator

.PHONY: build run test clean

build:
	go build -o tron_eggs_sim

run:
	go run main.go

test:
	go test ./...

test-verbose:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f tron_eggs_sim
