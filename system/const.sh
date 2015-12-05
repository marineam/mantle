#!/bin/bash -e

export GOOS=linux
for GOARCH in amd64; do
    export GOARCH
    go tool cgo -godefs const_${GOOS}.go > const_${GOOS}_${GOARCH}.go
    echo const_${GOOS}_${GOARCH}.go
done
