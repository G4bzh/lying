#!/bin/bash

set -e

cd /etc/bind
wget -q "http://dnscfg:8053/$CLIENTID/config" -O named.conf
for zone in $(wget -q  "http://dnscfg:8053/$CLIENTID/config/zones" -O-);
do
  wget -q  "http://dnscfg:8053/$CLIENTID/config/zone/$zone" -O $zone.txt;
done

exec "$@"
