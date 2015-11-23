#!/bin/bash

FILES="*.go"
TEST="flt05"
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
