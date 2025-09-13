#!/bin/bash
set -euo pipefail
sudo apt update
echo "Installing Dependencies"
sudo apt install -y podman podman-compose golang
echo "Setting up dir structure"
chmod +x deployment/dir_setup.sh
./deployment/dir_setup.sh
echo "Starting Containers"
cd deployment
podman-compose up -d
cd ..
sleep 10
echo "Starting the app"
go run ./cmd/main.go
echo "App stopped"
