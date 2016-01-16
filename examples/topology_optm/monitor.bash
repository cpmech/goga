#!/bin/bash

GOGA="$HOME/10.go/pkg/linux_amd64/github.com/cpmech/goga.a"
FILES="*.go *.json *.msh *.sim"

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
    #go run genmsh.go
    go run femsim.go topology.go
    #go run setandrunfem.go
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
