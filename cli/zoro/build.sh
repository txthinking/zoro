#!/bin/bash

if [ $# -ne 1 ]; then
    echo "./build.sh version"
    exit
fi

mkdir _

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o _/zoro_darwin_arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o _/zoro_darwin_amd64
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags="-w -s" -o _/zoro_freebsd_386
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags="-w -s" -o _/zoro_freebsd_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-w -s" -o _/zoro_linux_386
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o _/zoro_linux_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o _/zoro_linux_arm64
CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build -ldflags="-w -s" -o _/zoro_netbsd_386
CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build -ldflags="-w -s" -o _/zoro_netbsd_amd64
CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -ldflags="-w -s" -o _/zoro_openbsd_386
CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -ldflags="-w -s" -o _/zoro_openbsd_amd64
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -ldflags="-w -s" -o _/zoro_openbsd_arm64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o _/zoro_windows_amd64.exe
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-w -s" -o _/zoro_windows_386.exe

nami release github.com/txthinking/zoro $1 _

rm -rf _
