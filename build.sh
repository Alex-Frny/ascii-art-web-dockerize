#!/bin/sh

set -eu

IMAGE_NAME="ascii-art-web:1.0"
CONTAINER_NAME="ascii-art-web-container"

echo "Construction de l'image Docker..."
docker image build \
  -f Dockerfile \
  -t "$IMAGE_NAME" \
  .

if docker container inspect "$CONTAINER_NAME" >/dev/null 2>&1
then
    echo "Suppression de l'ancien conteneur..."
    docker container rm -f "$CONTAINER_NAME"
fi

echo "Création du conteneur..."
docker container run \
  --detach \
  --name "$CONTAINER_NAME" \
  --label project="ascii-art-web-dockerize" \
  --label author="Alexis Fourny" \
  --publish 8080:8080 \
  "$IMAGE_NAME"

echo "Conteneur actif :"
docker container ps \
  --filter "name=$CONTAINER_NAME"