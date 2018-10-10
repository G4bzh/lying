#!/bin/bash

echo "Provisioning client virtual machine..."


sudo su -
wget -O - 'https://dl.bintray.com/rabbitmq/Keys/rabbitmq-release-signing-key.asc' | sudo apt-key add -
echo "deb http://dl.bintray.com/rabbitmq/debian stretch erlang" > /etc/apt/sources.list.d/bintray.erlang.list
apt-get update -y
apt-get install erlang-nox -y
echo "deb https://dl.bintray.com/rabbitmq/debian stretch main" > /etc/apt/sources.list.d/bintray.rabbitmq.list
apt-get update -y
apt-get install rabbitmq-server -y
service rabbitmq-server start
