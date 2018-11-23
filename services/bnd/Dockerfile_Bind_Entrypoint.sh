#!/bin/bash

# Exit on error
set -e

cd /etc/bind
wget -q "http://dnscfg:8053/v1/private/$CLIENTID/config" -O named.conf
for zone in $(wget -q  "http://dnscfg:8053/v1/private/$CLIENTID/config/zones" -O-);
do
  wget -q  "http://dnscfg:8053/v1/private/$CLIENTID/config/zone/$zone" -O $zone.txt;
done

exec "$@"
