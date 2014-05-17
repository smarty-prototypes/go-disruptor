#!/bin/bash
set -e

GOOS=linux GOARCH=arm GOARM=7 go build -o go-disruptor
adb push go-disruptor /data/local/tmp
adb shell "cd /data/local/tmp; chmod 755 go-disruptor; ./go-disruptor"
