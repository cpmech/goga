#!/bin/bash

FILES="*.go"
TEST="flt04"

refresh(){
    echo
    echo
    go test -test.run="$TEST"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
