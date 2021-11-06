#/bin/bash

set -exuo pipefail

arr=$(find . | grep '99_hw$' | grep -v '4/99_hw$' | grep -v '99_hw/code' )
for i in $arr; do golangci-lint -c .golangci.yml run $i/...;done

cd 4/99_hw/taskbot
golangci-lint -c ../../../.golangci.yml run --modules-download-mode=vendor ./...