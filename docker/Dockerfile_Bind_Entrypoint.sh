#!/bin/bash

set -e

cd /etc/bind
/etc/bind/toZone -id "$CLIENTID" -url "$MONGOURL"

exec "$@"
