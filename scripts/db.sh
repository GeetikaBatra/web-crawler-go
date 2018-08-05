#!/bin/bash

# fail if smth fails
# the whole env will be running if test suite fails so you can debug
set -e

set -x

here=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)
DATABASE="crawl"
USERNAME="crawl"
export PGPASSWORD="crawl"

POSTGRES_CONTAINER_NAME="crawl-postgres"
DB_CONTAINER_IP=$(docker inspect --format '{{.NetworkSettings.IPAddress}}' ${POSTGRES_CONTAINER_NAME})

# TODO: this is duplicating code with server's runtest, we should refactor
echo "Waiting for postgres to fully initialize"
set +x
for i in {1..10}; do
  retcode=`curl http://${DB_CONTAINER_IP}:5432 &>/dev/null || echo $?`
  if test "$retcode" == "52"; then
    break
  fi;
  sleep 1
done;
psql -h $DB_CONTAINER_IP -U $USERNAME $DATABASE << EOF
DROP TABLE IF EXISTS test;
CREATE TABLE test(name VARCHAR (50) PRIMARY KEY);
EOF
set -x