#!/bin/bash

echo "Provisioning client virtual machine..."


sudo su -
apt-get update -y
apt-get install dnsutils -y
wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
tar -xvzf go1.11.1.linux-amd64.tar.gz -C /usr/local/
echo "export PATH=$PATH:/usr/local/go/bin" > /etc/profile
apt-get install -y git
