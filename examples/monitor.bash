#!/bin/bash

#FILE="rel-simple-beam-form.go"
FILE="reliability-problem1.go"

while true; do
    inotifywait -q -e modify $FILE
    go run $FILE
done
