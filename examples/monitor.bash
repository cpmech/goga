#!/bin/bash

FILE="rel-simple-beam-form.go"

while true; do
    inotifywait -q -e modify $FILE
    go run $FILE
done
