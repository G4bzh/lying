#!/bin/bash

echo "Provisioning bind virtual machine..."


sudo su -
apt-get update -y
apt-get install bind9 dnsutils -y
