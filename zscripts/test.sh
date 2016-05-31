#!/bin/sh
cd $GOPATH/src/github.com/jimmy-go/pgwp
go test -cover -coverprofile=coverage.out

if [ "$1" == "html" ]; then
    go tool cover -html=coverage.out
fi
