#/bin/bash

arr=$(find . | grep '99_hw$')
for i in $arr; do golangci-lint -c .golangci.yml run $i/...;done
