#!/bin/bash

echo "Provisioning bind virtual machine..."


sudo su -
apt-get update -y
apt-get install bind9 dnsutils -y
mkdir /var/log/bind
chown bind:bind /var/log/bind
echo 'include "/bind/named.conf.log";' >> /etc/bind/named.conf
