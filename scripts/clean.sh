#!/bin/bash
## DeGOps: 0.0.4
set -o errexit
set -o nounset

rm -rf vendor
touch coverage.out && rm coverage.out
touch coverage.html && rm coverage.html
