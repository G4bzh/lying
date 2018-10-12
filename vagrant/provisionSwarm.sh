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

# Overlay Network
# docker network create --subnet 1.0.0.0/8 --driver overlay saas

# Consul Server
# docker service create --name consul --network saas --mode global --env CONSUL_BIND_INTERFACE=eth0 --env CONSUL_CLIENT_INTERFACE=eth0 consul agent -server -bootstrap-expect=3  -retry-join=tasks.consul
#
# Test : curl http://consul:8500/v1/agent/members?pretty (from within a container, docker ps thent docker exec -it <ID> /bin/sh)
