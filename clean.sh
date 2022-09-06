#!/bin/bash

set -e

rm -rf ./prelease/prelease

rm -rf ./server/__debug_bin
rm -rf ./server/server
rm -rf ./server/tmp*
# rm -rf ./server/data

if [[ $1 == 'all' ]] 
then

    rm -rf ./server/build

else

    rm -rf ./server/build/linux64
    
fi

rm -rf ./server/build/*.gz
