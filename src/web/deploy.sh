#!/bin/sh

echo "Deploying forklol-web $1"
sed -i -E "s/VERSION: '\"v(.*)\"'/VERSION: '\"$1\"'/g" ./config/dev.env.js
sed -i -E "s/VERSION: '\"v(.*)\"'/VERSION: '\"$1\"'/g" ./config/prod.env.js

echo "Running node build"
VERSION=$1 npm run build

echo "Building docker container"
VERSION=$1 docker build -t forklol-web:$1 .

echo "Tagging container for GCE"
docker tag forklol-web:$1 eu.gcr.io/splendid-sled-172714/forklol-web:$1

echo "Pushing container to GCE"
gcloud docker -- push eu.gcr.io/splendid-sled-172714/forklol-web:$1

echo "Replacing kube deployment"
sed -i -E "s/value: v(.*) #VERSION/value: $1 #VERSION/g" ./kube.deploy.yaml
sed -i -E "s/forklol-web:v(.*)/forklol-web:$1/g" ./kube.deploy.yaml
kubectl replace -f ./kube.deploy.yaml