#!/bin/bash

GP1="${GOPATH%:*}"
GP2="${GOPATH#*:}"

GP=$GP2
if [[ -z "${GP// }" ]]; then
    GP=$GP1
fi

GOGA="$GP/pkg/linux_amd64/github.com/cpmech/goga.a"

FILES="*.go"

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
    go run two-obj-ct.go
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
