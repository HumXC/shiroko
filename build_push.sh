#!/usr/bin/env sh
export GOOS=linux
export GOARCH=arm64
go build -o ./build/shiroko ./cmd/shiroko
adb push build/shiroko /data/local/tmp/ >/dev/null
adb shell "/data/local/tmp/shiroko $@"
