#!/bin/bash

FILES="*.go"

refresh(){
    echo
    echo
    go test -test.run="functions"
    #go test -test.run="flt05"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
