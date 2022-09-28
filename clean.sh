#!/bin/bash

set -e

rm -rf ./server/__debug_bin
rm -rf ./server/server
rm -rf ./server/tmp*
# rm -rf ./server/data
rm -rf ./server/api/db/data

rm -rf ./data/out

if [[ $1 == 'all' ]] 
then

    rm -rf ./server/build

else

    rm -rf ./server/build/linux64/server*
    rm -rf ./server/build/linux64/process
    rm -rf ./server/build/linux64/tmp-locker
    
fi

rm -rf ./server/build/*.gz
