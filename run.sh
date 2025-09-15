#!/bin/bash
set -euo pipefail
cd ./public/
git pull
cd ../deployment/
podman-compose up -d
cd ../cmd/
go run ./main.go
