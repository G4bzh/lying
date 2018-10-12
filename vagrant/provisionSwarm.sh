#!/bin/bash

echo "Provisioning client virtual machine..."


sudo su -
apt-get update -y
apt-get install apt-transport-https ca-certificates curl gnupg2 software-properties-common -y
curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
apt-get update -y
apt-get install docker-ce -y

# Swarm01
# docker swarm init --advertise-addr 192.168.10.1

# Swarm 02/03
# docker swarm join --token SWMTKN-1-4b0yzdn30rx7qopbor2j0frbdx110a9kxfurtpqtp4modi80so-25gg795z2y85tj47j7hv56w90 192.168.10.1:2377

# docker network create --subnet 10.0.0.0/8 --driver overlay saas
