#!/bin/bash

echo "Provisioning bind virtual machine..."


sudo su -
apt-get update -y
apt-get install bind9 dnsutils -y
mkdir /var/log/bind
chown bind:bind /var/log/bind
echo 'include "/bind/named.conf.log";' >> /etc/bind/named.conf
echo 'include "/bind/named.conf.zones";' >> /etc/bind/named.conf.local
wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
tar -xvzf go1.11.1.linux-amd64.tar.gz -C /usr/local/
echo "export PATH=$PATH:/usr/local/go/bin" > /etc/profile
apt-get install -y git
