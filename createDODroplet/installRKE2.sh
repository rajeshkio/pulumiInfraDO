#!/bin/bash
echo "Updating ubuntu and installing curl" > /tmp/upgrade.log
apt update && apt install curl -y

echo "Installing rke2 server" > /tmp/upgrade.log 
curl -sfL https://get.rke2.io | sh -
systemctl enable rke2-server.service
systemctl start rke2-server.service

if [[ $? -eq 0 ]]; then 
  echo "RKE2 server installed successfully" > /tmp/upgrade.log
fi
