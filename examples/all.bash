#!/bin/bash

set -e

examples="\
01-one-obj \
02-two-obj \
03-two-obj-ct \
04-three-obj \
05-many-obj \
06-truss \
07-eed 
"

for ex in $examples; do
    echo
    echo
    echo ">>> running $ex <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
    cd $ex
    ./all.bash
    cd ..
done
