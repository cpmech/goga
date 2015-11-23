#!/bin/bash

FILES="*.go"
TEST="flt01"
#TEST="sort01"

refresh(){
    echo
    echo
    go test -test.run="$TEST"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
