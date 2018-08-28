#!/bin/bash

# fail if smth fails
# the whole env will be running if test suite fails so you can debug
set -e

set -x

here=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)

DB_CONTAINER_IP=cassandra-docker
# TODO: this is duplicating code with server's runtest, we should refactor
#!/bin/bash

status=$(nc -z cassandra-docker 9042; echo $?)
echo $status

while [ $status != 0 ]
do
  sleep 3s
  status=$(nc -z cassandra-docker 9042; echo $?)
  echo $status
done

exec ./bin/gremlin-server.sh ./conf/gremlin-server/http-gremlin-server.yaml