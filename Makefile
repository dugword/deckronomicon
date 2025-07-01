# Makefile for Tron Eggs Simulator

.PHONY: build run test clean

run:
	go run ./

test:
	go test ./...

test-verbose:
	go test -v ./...

lint:
	golangci-lint run
