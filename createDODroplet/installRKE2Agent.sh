#!/bin/bash

apt update && apt install curl -y
curl -sfL https://get.rke2.io | INSTALL_RKE2_TYPE="agent" sh -

mkdir -p /etc/rancher/rke2/
server: https://MASTER_IP:9345 > /etc/rancher/rke2/config.yaml
token: NODE_TOKEN >> /etc/rancher/rke2/config.yaml

systemctl start rke2-agent.service