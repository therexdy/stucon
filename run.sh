#!/bin/bash
if [ "$1" = "--stop" ]; then
  podman-compose down

  cd ./deployment/
  podman-compose down
  exit
fi

cd ./public/
git pull
cd ..

git pull

podman network create deployment_stucon --subnet=192.168.200.0/24

cd ./deployment/
podman-compose down
podman-compose up -d

cd ../
podman-compose down

if [ "$1" = "--build" ]; then
  podman-compose up --build -d
else
  podman-compose up -d
fi
