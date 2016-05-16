#!/bin/bash

GOGA="$HOME/10.go/pkg/linux_amd64/github.com/cpmech/goga.a"
FILES="*.go"

if [ -f $GOGA ]; then
   FILES="$FILES $GOGA"
fi

echo
echo "monitoring:"
echo $FILES
echo
echo

refresh(){
    echo
    echo
    #go run one-obj.go
    #go run one-obj-prob9.go
    go run one-obj-prob9-dbg.go
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
