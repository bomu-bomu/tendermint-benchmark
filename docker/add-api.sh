#!/bin/sh

ip="$(ifconfig | grep -A 1 'eth0' | tail -1 | cut -d ':' -f 2 | cut -d ' ' -f 1)"
curl "http://${SERVER_HOSTNAME}:${SERVER_PORT}/ip/add/${ip}"
