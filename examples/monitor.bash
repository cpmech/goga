#!/bin/bash

FILES="tsp-simple.go tsp-simple.json"
GOFILE="tsp-simple.go"

while true; do
    inotifywait -q -e modify $FILES
    go run $GOFILE
done
