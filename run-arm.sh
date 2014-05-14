#!/bin/bash
set -e

GOOS=linux GOARCH=arm GOARM=7 go build
~/Downloads/android/adt-bundle-mac-x86_64-20140321/sdk/platform-tools/adb push go-disruptor /data/local/tmp
~/Downloads/android/adt-bundle-mac-x86_64-20140321/sdk/platform-tools/adb shell "chmod 755 /data/local/tmp/go-disruptor"
~/Downloads/android/adt-bundle-mac-x86_64-20140321/sdk/platform-tools/adb shell "/data/local/tmp/go-disruptor"
