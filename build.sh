#!/bin/bash

set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

cd ./server

GOARCH=amd64
LDFLAGS="-s -w"
TM=`date +%F@%T@%Z`
OUT=server\($TM\)

# For Docker, one build below for linux64 is enough.
OUTPATH_LINUX=./build/linux64/
mkdir -p $OUTPATH_LINUX
CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH_LINUX
echo "${G}server(linux64) built${W}"

#######################################################################################

if [[ $1 == 'release' || $1 == 'rel' ]]
then

    RELEASE_NAME=wisite-api\($TM\).tar.gz 
    cd ./build
    echo $RELEASE_NAME
    tar -czvf ./$RELEASE_NAME --exclude='./linux64/data' ./linux64

fi