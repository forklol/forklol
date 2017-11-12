#!/bin/sh

DIR="`dirname \"$0\"`"
DIR="`( cd \"$DIR\" && pwd )`"

cd $DIR

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o forklol-api .

docker image rm -f fl-backend:latest | true
docker build -t fl-backend:latest .
