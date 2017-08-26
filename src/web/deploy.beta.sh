#!/bin/sh

echo "Deploying forklol-web $1b"
sed -i -E "s/VERSION: '\"v(.*)\"'/VERSION: '\"$1b\"'/g" ./config/dev.env.js
sed -i -E "s/VERSION: '\"v(.*)\"'/VERSION: '\"$1b\"'/g" ./config/prod.env.js

echo "Running node build"
VERSION=$1b npm run build

echo "Building docker container"
VERSION=$1b docker build -t forklol-web:$1b .

echo "Tagging container for GCE"
docker tag forklol-web:$1b eu.gcr.io/splendid-sled-172714/forklol-web:$1b

echo "Pushing container to GCE"
gcloud docker -- push eu.gcr.io/splendid-sled-172714/forklol-web:$1b

echo "Replacing kube deployment"
sed -i -E "s/value: v(.*) #VERSION/value: $1b #VERSION/g" ./kube.deploy.beta.yaml
sed -i -E "s/forklol-web:v(.*)/forklol-web:$1b/g" ./kube.deploy.beta.yaml
kubectl replace -f ./kube.deploy.beta.yaml