#!/bin/bash

echo "Provisioning client virtual machine..."


sudo su -
apt-get update -y
apt-get install apt-transport-https ca-certificates curl gnupg2 software-properties-common -y
curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
apt-get update -y
apt-get install docker-ce dnsutils -y
mkdir -p ~/.ssh 2>/dev/null
cp /vagrant/id_rsa* ~/.ssh
chmod 600 ~/.ssh/id_rsa
cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys
if [ "$HOSTNAME" == "swarm01" ]; then
  docker swarm init --advertise-addr 192.168.10.1
  docker swarm join-token -q manager > /vagrant/swarm_token
  docker node update swarm01 --label-add mongodb=primary
  ping -c 3 192.168.10.2
  ssh-keyscan -H 192.168.10.2 >> ~/.ssh/known_hosts
  ssh 192.168.10.2 "docker swarm join --token `echo -n $(head /vagrant/swarm_token)` 192.168.10.1:2377"
  ping -c 3 192.168.10.3
  ssh-keyscan -H 192.168.10.3 >> ~/.ssh/known_hosts
  ssh 192.168.10.3 "docker swarm join --token `echo -n $(head /vagrant/swarm_token)` 192.168.10.1:2377"
fi
