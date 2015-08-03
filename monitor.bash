#!/bin/bash

FILE="t_evolver_test.go"
#FILE="evolver.go"
TEST="evo03"

refresh(){
    go test -test.run="$TEST"
}

while true; do
    inotifywait -q -e modify $FILE
    refresh
done
