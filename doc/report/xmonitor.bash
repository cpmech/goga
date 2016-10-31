#!/bin/bash

FILES="*.tex *.sty"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    make
done
