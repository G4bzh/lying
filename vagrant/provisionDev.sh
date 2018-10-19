#!/bin/bash

echo "Provisioning bind virtual machine..."


sudo su -
apk update
apk add go
apk add mongodb
apk add mongodb-tools
apk add git
apk add musl-dev
mkdir -p /data/db
exit
go get github.com/globalsign/mgo
