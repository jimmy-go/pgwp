#!/bin/bash
## DeGOps: 0.0.4
set -o errexit
set -o nounset

go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
