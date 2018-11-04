#!/bin/bash

set -e

cd /etc/bind
wget -q "http://dnscfg:8080/$CLIENTID/config" > named.conf
for zone in $(wget -q  "http://dnscfg:8080/foo@bar.com/zones" -O-);
do
  wget -q  "http://dnscfg:8080/$CLIENTID/zone/$zone" -O- > $zones.txt;
done

exec "$@"
