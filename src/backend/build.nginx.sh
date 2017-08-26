#!/bin/sh

echo "Building docker container"
docker build -t forklol-api-nginx:$1 . --file ./Dockerfile.nginx

echo "Tagging container for GCE"
docker tag forklol-api-nginx:$1 eu.gcr.io/splendid-sled-172714/forklol-api-nginx:$1

echo "Pushing container to GCE"
gcloud docker -- push eu.gcr.io/splendid-sled-172714/forklol-api-nginx:$1