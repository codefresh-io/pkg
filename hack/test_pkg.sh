#!/bin/sh

set -e
rm -f coverage.txt || true
covs=$(find . -mindepth 2 -maxdepth 4 -type f -name 'coverage.txt')

for d in $covs; do
    echo $d
    if [ -f $d ]; then
        cat $d >> coverage.txt
        rm $d
    fi
done

