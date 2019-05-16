#!/bin/bash

if [ ! -d binary ]
then
    mkdir binary
fi

if [ -f binary.tgz ]
then
    rm binary.tgz
fi

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o binary/mr2 .
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o binary/mr2_linux_386 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o binary/mr2_linux_arm64 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o binary/mr2_linux_arm7 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o binary/mr2_linux_arm6 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o binary/mr2_linux_arm5 .
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o binary/mr2_linux_mips .
CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o binary/mr2_linux_mips_sf .
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -o binary/mr2_linux_mipsle .
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -o binary/mr2_linux_mipsle_sf .
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -o binary/mr2_linux_mips64 .
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -o binary/mr2_linux_mips64le .
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build -o binary/mr2_linux_ppc64 .
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -o binary/mr2_linux_ppc64le .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o binary/mr2_darwin_amd64 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o binary/mr2_windows_amd64.exe .
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o binary/mr2_windows_386.exe .

tar zcf binary.tgz binary
