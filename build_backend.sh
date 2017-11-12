#!/bin/sh

DIR="`dirname \"$0\"`"
DIR="`( cd \"$DIR\" && pwd )`"

cd $DIR
cd ./src/backend

go clean
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o forklol-api -v .
