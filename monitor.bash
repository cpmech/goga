#!/bin/bash

FILES="evolver.go island.go t_evolver_test.go"
TEST="evo03"

refresh(){
    go test -test.run="$TEST"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
