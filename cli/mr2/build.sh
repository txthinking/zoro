#!/bin/bash
GOOS=darwin GOARCH=386 go build -o mr2_darwin_386
GOOS=darwin GOARCH=amd64 go build -o mr2_darwin_amd64
GOOS=freebsd GOARCH=386 go build -o mr2_freebsd_386
GOOS=freebsd GOARCH=amd64 go build -o mr2_freebsd_amd64
GOOS=linux GOARCH=386 go build -o mr2_linux_386
GOOS=linux GOARCH=amd64 go build -o mr2_linux_amd64
GOOS=linux GOARCH=arm64 go build -o mr2_linux_arm64
GOOS=netbsd GOARCH=386 go build -o mr2_netbsd_386
GOOS=netbsd GOARCH=amd64 go build -o mr2_netbsd_amd64
GOOS=openbsd GOARCH=386 go build -o mr2_openbsd_386
GOOS=openbsd GOARCH=amd64 go build -o mr2_openbsd_amd64
GOOS=openbsd GOARCH=arm64 go build -o mr2_openbsd_arm64

GOOS=windows GOARCH=amd64 go build -o mr2_windows_amd64.exe .
GOOS=windows GOARCH=386 go build -o mr2_windows_386.exe .

rm mr2.tgz
tar czf mr2.tgz mr2_*
rm -rf mr2_*
