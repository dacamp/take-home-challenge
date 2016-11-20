#!/bin/bash

for i in {1..10}; do
    curl -m 6  -X POST -d "$i"  http://localhost:1234/counter/foobar &
    curl -m 6 -X POST -d "$i"  http://localhost:1235/counter/foobar &
    curl -m 6 -X POST -d "$i"  http://localhost:1236/counter/foobar &
done


for node in 1234 1235 1236; do
    echo "NODE: ${node}"
    echo -n "VALUE: " && curl -m 6 http://localhost:${node}/counter/foobar/value
    echo -n "CONSISTENT VALUE: " && curl -m 6 http://localhost:${node}/counter/foobar/consistent_value
    echo
done
