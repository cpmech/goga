#!/bin/bash

FILES="*.go"

echo
echo "monitoring:"
echo $FILES
echo
echo

refresh(){
    echo
    echo
    go install
    go test -test.run="flt05"
    go test -test.run="flt06"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
