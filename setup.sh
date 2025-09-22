#!/bin/bash
set -euo pipefail
sudo apt update
echo "Installing Dependencies"
sudo apt install -y podman podman-compose 
echo "Setting up dir structure"
chmod +x deployment/dir_setup.sh
./deployment/dir_setup.sh
echo "Cloning Repo:"
mkdir -p public
cd ./public
git clone https://github.com/Chiranthcs6/Website.git
mv Website/* Website/.* ./
rmdir Website
cd ..
echo "Setup Complete, now do bash run.sh to start the server."
