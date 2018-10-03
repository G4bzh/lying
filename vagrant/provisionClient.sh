#!/bin/bash

echo "Provisioning client virtual machine..."


sudo su -
apt-get update -y
apt-get install dnsutils -y
