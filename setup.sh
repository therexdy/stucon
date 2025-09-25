#!/bin/bash
set -euo pipefail
# update
echo "Installing Dependencies"
# install -y podman podman-compose 
echo "Setting up dir structure"
chmod +x deployment/dir_setup.sh
./deployment/dir_setup.sh
echo "Cloning Repo:"
mkdir -p public
cd ./public
git clone https://github.com/Chiranthcs6/stucon_Connect.git
repo_name="stucon_Connect"
mv $repo_name/* $repo_name/.* ./
rmdir $repo_name
cd ..
echo "Setup Complete, now do bash run.sh to start the server."
