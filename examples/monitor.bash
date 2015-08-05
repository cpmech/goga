#!/bin/bash

FILE="reliability-simple.go"

while true; do
    inotifywait -q -e modify $FILE
    go run $FILE
done
