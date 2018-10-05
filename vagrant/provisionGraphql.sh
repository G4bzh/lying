#!/bin/bash

echo "Provisioning client virtual machine..."


sudo su -
apt-get update -y
apt-get install -y build-essential cur
curl -sL https://deb.nodesource.com/setup_8.x |  bash -
echo "Package: *" > /etc/apt/preferences.d/nodejs
echo "Pin: origin deb.nodesource.com" >> /etc/apt/preferences.d/nodejs
echo "Pin-Priority: 1002" >> /etc/apt/preferences.d/nodejs
apt-get update -y
apt-get install -y nodejs
mkdir /app
cd /app
npm init -y
npm install express express-graphql graphql --save
