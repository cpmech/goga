#!/bin/bash

FILES="reliability-simple tsp-simple"

for f in $FILES; do
    echo
    echo
    echo "[1;33m>>>>>>>>>>>>>>>>>>>>>>>>>>> running $f <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    echo
    go build -o /tmp/goga/$f "$f".go && /tmp/goga/$f
done
