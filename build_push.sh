#!/usr/bin/env sh
export GOOS=linux
export GOARCH=arm64
go build -o shiroko
adb push shiroko /data/local/tmp/
adb shell "/data/local/tmp/shiroko $@"
