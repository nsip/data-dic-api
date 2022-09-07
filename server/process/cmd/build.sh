#!/bin/bash

set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

GOARCH=amd64
LDFLAGS="-s -w"
OUT=process

CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT

# move it to server executable directory
SERVER_PATH=../../build/linux64/
mv $OUT $SERVER_PATH
echo "${G}process(linux64) built, and moved to $SERVER_PATH ${W}"