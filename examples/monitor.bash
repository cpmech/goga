#!/bin/bash

#FILE="rel-simple-beam-form.go"
FILE="rel-2d-simple.go"

while true; do
    inotifywait -q -e modify $FILE
    go run $FILE
done
