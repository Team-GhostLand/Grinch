#!/bin/bash

# Run this command from the main repo dir (ie. as ./scripts/build.sh) on an x64 Linux machine.

rm -r "./bin/"; #If you didn't have a ./bin/ directory yet, ignore any errors this might throw. The script will carry on.
go build -o "./bin/grinch";
GOOS=windows go build -o "./bin/grinch.exe";
GOOS=windows GOARCH=arm go build -o "./bin/grinch-arm.exe"; #No need to specify arm64; Windows ARM is implicitly 64bit
GOOS=darwin GOARCH=arm64 go build -o "./bin/grinch-mac";
GOARCH=arm64 go build -o "./bin/grinch-arm";
GOARCH=riscv64 go build -o "./bin/grinch-r5";
zip -9jq "bin/linux.zip" "bin/grinch" "bin/grinch-arm" "bin/grinch-r5";
ln -rs ./scripts/make_serverpack.sh ./bin/make_serverpack.sh