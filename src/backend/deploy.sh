#!/bin/sh

echo "Building binary for $1"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o forklol-api .

echo "Building docker container"
docker build -t forklol-api:$1 .

echo "Tagging container for GCE"
docker tag forklol-api:$1 eu.gcr.io/splendid-sled-172714/forklol-api:$1

echo "Pushing container to GCE"
gcloud docker -- push eu.gcr.io/splendid-sled-172714/forklol-api:$1

echo "Replacing kube deployment"
sed -i -E "s/forklol-api:v(.*)/forklol-api:$1/g" ./kube.deploy.proxy.yaml
kubectl replace -f ./kube.deploy.proxy.yaml