#!/bin/bash
#Expect to run this from the root of the project as "./build/buildImage.sh"

DOCKER_TAG="$(date +%s)"
echo "DOCKER_TAG: $DOCKER_TAG"

GIT_COMMIT="$(git rev-parse HEAD)"
echo "GIT_COMMIT hash is $GIT_COMMIT"

#Set config file reference
CONFIG_FILE="deploy/portfolio-service-conf_standalone.json"

docker build --build-arg "GIT_COMMIT=GIT_COMMIT" --build-arg "CONFIG_FILE=$CONFIG_FILE" --tag "portfolio-service:$DOCKER_TAG" --tag "portfolio-service:latest" -f build/Dockerfile .