#!/bin/bash
set -euo pipefail
cd ./public/
git pull
cd ../cmd/
go run ./main.go
