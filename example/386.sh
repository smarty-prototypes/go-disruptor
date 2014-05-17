#!/bin/bash
set -e

GOOS=linux GOARCH=386 go build -o go-disruptor
