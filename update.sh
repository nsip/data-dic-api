#!/bin/bash

set -e

rm -rf ./server/docs

cd ./server
./swagger/swag init
cd -

if [[ $1 == 'all' ]]
then

rm -f go.sum
go get -u ./...

fi