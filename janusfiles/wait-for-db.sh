#!/bin/bash

# fail if smth fails
# the whole env will be running if test suite fails so you can debug
set -e

set -x

here=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)

DB_CONTAINER_IP=cassandra-docker
# TODO: this is duplicating code with server's runtest, we should refactor
echo "Waiting for cassandra to fully initialize"
set +x
for i in {1..50}; do
  retcode=`curl ${DB_CONTAINER_IP}:9042 &>/dev/null || echo $?`
  if test "$retcode" == "52"; then
    break
  fi;
  sleep 1
done;

set -x