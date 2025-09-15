#!/bin/bash
set -euo pipefail
cd ./public/
git pull
cd ../deployment/
podman-compose up -d
cd ../
podman-compose up -d
