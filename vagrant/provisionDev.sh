#!/bin/bash

echo "Provisioning bind virtual machine..."


sudo su -
apt-get update -y
wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
tar -xvzf go1.11.1.linux-amd64.tar.gz -C /usr/local/
echo "export PATH=$PATH:/usr/local/go/bin" > /etc/profile
apt-get install -y git
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
# MongoDB 3.6 for mgo driver (4.0 not supported yet)
echo "deb http://repo.mongodb.org/apt/debian stretch/mongodb-org/3.6 main" | tee /etc/apt/sources.list.d/mongodb-org-3.6.list
apt-get update -y
apt-get install -y mongodb-org
