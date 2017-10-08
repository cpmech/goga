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

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go run simple02.go
    #go run cross-in-tray.go
    go run cross-in-tray-stat.go
done
