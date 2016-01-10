#!/bin/bash

FILES="*.go"

refresh(){
    echo
    echo
    go run simpleoptm.go
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
