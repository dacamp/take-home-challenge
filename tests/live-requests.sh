#!/bin/bash
# NOTE: not even close to 'web scale'
#
# Super jank integration hack tests:
#  * container on :1234 can receive configs
#  * POST requests to containers on ports :1235 & :1236 can receive
#    counter values
#  * those values get replicated out to all nodes, include :1234


if [[ ! -z "$1" ]]; then
    echo "UPDATING CONFIG"
    curl -i \
	 -H "Accept: application/json" \
	 -H "Content-Type:application/json" \
	 -X POST --data '{ "actors": ["172.17.0.3", "172.17.0.2", "172.17.0.4"] }' http://localhost:1234/config
fi

for i in {1..10}; do
    #curl -m 6  -X POST -d "$i"  http://localhost:1234/counter/foobar &
    curl -m 6 -X POST -d "$i"  http://localhost:1235/counter/foobar &
    curl -m 6 -X POST -d "$i"  http://localhost:1236/counter/foobar &
done


for node in 1234 1235 1236; do
    echo "NODE: ${node}"
    echo -n "VALUE: " && curl -m 6 http://localhost:${node}/counter/foobar/value
    echo -n "CONSISTENT VALUE: " && curl -m 6 http://localhost:${node}/counter/foobar/consistent_value
    echo
done
