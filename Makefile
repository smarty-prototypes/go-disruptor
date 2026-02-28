#!/usr/bin/make -f

test:
	go test -timeout=1s -short -race -covermode=atomic ./...

test.long:
	go test -run TestEndToEnd -timeout=120s -race -covermode=atomic -v ./...

benchmark:
	go test -bench=. -benchmem

compile:
	go build ./...

build: test compile

.PHONY: test test.long compile build benchmark
