#!/usr/bin/make -f

test:
	go test -timeout=1s -race -covermode=atomic ./...

benchmark:
	go test -bench=. -benchmem

compile:
	go build ./...

build: test compile

.PHONY: test compile build benchmark
