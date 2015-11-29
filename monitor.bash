#!/bin/bash

FILES="*.go"

refresh(){
    echo
    echo
    go test -test.run="flt05"
    go test -test.run="flt06"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
