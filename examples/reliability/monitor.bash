#!/bin/bash

GOGA="$GOPATH/pkg/linux_amd64/github.com/cpmech/goga.a"
FILES="*.go *.json *.py"

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
    #go run defs.go fesim.go ReliabFORM.go simple.go  simple
    go run defs.go fesim.go ReliabFORM.go simple.go  frame2d
    #python draw-frame2d.py
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
