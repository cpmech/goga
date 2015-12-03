#!/bin/bash

FILES="*.go"

refresh(){
    echo
    echo
    go test -test.run="functions"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
