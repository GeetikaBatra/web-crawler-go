#!/bin/bash

# fail if smth fails
# the whole env will be running if test suite fails so you can debug
set -e

set -x

JANUS_IP="janus"

while :
do
  sleep 3s
  status=$(curl ${JANUS_IP}:8182; echo $?)
  echo $status
  if test "$status" != "7"; then
    break
  fi
done
echo "Final"
exec ./main https://monzo.com/