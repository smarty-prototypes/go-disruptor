#!/bin/bash
set -e

export PATH="$PATH:~/Downloads/android/adt-bundle-mac-x86_64-20140321/sdk/platform-tools"

case "$1" in
	copy)
		GOOS=linux GOARCH=arm GOARM=7 go build -o go-disruptor
		adb push go-disruptor /data/local/tmp
	;;
esac

adb shell "cd /data/local/tmp; chmod 755 go-disruptor; ./go-disruptor"
