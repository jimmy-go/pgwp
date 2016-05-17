#!/bin/sh
. $GOPATH/src/github.com/jimmy-go/pgwp/zscripts/envs.sh
cd $GOPATH/src/github.com/jimmy-go/pgwp/examples/basic

go build -o $GOBIN/pgwp_basic

gcvis $GOBIN/pgwp_basic \
-host=$PG_HOST \
-database=$PG_DATABASE \
-u=$PG_USERNAME \
-p=$PG_PASSWORD
