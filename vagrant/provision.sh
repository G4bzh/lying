#!/bin/bash

echo "Provisioning virtual machine..."


sudo su -
apt-get update -y
apt-get install bind9 -y 
