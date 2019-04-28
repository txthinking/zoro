#!/bin/bash

if [ ! -d binary ]
then
    mkdir binary
fi
exit;

GOOS=linux GOARCH=amd64 go build -o binary/mr2 .
GOOS=linux GOARCH=386 go build -o binary/mr2_linux_386 .
GOOS=linux GOARCH=arm64 go build -o binary/mr2_linux_arm64 .
GOOS=linux GOARCH=arm GOARM=7 go build -o binary/mr2_linux_arm7 .
GOOS=linux GOARCH=arm GOARM=6 go build -o binary/mr2_linux_arm6 .
GOOS=linux GOARCH=arm GOARM=5 go build -o binary/mr2_linux_arm5 .
GOOS=linux GOARCH=mips go build -o binary/mr2_linux_mips .
GOOS=linux GOARCH=mipsle go build -o binary/mr2_linux_mipsle .
GOOS=linux GOARCH=mips64 go build -o binary/mr2_linux_mips64 .
GOOS=linux GOARCH=mips64le go build -o binary/mr2_linux_mips64le .
GOOS=linux GOARCH=ppc64 go build -o binary/mr2_linux_ppc64 .
GOOS=linux GOARCH=ppc64le go build -o binary/mr2_linux_ppc64le .
GOOS=darwin GOARCH=amd64 go build -o binary/mr2_darwin_amd64 .
GOOS=windows GOARCH=amd64 go build -o binary/mr2_windows_amd64.exe .
GOOS=windows GOARCH=386 go build -o binary/mr2_windows_386.exe .

tar zcf binary.tgz binary
