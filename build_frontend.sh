#!/bin/sh

DIR="`dirname \"$0\"`"
DIR="`( cd \"$DIR\" && pwd )`"

cd $DIR
cd ./src/web

sed -i -E "s/VERSION: '\"v(.*)\"'/VERSION: '\"$1\"'/g" ./config/dev.env.js
sed -i -E "s/VERSION: '\"v(.*)\"'/VERSION: '\"$1\"'/g" ./config/prod.env.js

docker run --rm -v $(pwd):/src -w /src node:6 bash -c "VERSION=$1 npm run build"
docker build -t fl-frontend:latest .
