#!/bin/bash

GP1="${GOPATH%:*}"
GP2="${GOPATH#*:}"

GP=$GP2
if [[ -z "${GP// }" ]]; then
    GP=$GP1
fi

GOGA="$GP/pkg/linux_amd64/github.com/cpmech/goga.a"

FILES="*.go *.json *.msh *.sim *.py"

if [ -f $GOGA ]; then
   FILES="$FILES $GOGA"
fi

echo
echo "monitoring:"
echo $FILES
echo
echo "with:"
echo "GP = $GP"
echo
echo

refresh(){
    echo
    echo
    echo
    echo
    echo
    #go run genmsh.go
    #go run drawmsh.go
    go run femsim.go reporting.go topology.go
    #go run setandrunfem.go
    #go run plotCPUtime.go
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
