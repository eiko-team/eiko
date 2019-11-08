#!/bin/bash
E_VAL=0

for f in $(find . -name '*.go');do
    TMP=$(gofmt -s $f | diff $f -);
    VAL=$?
    if [[ $VAL -ne 0 ]]; then
        E_VAL=$VAL
        echo $f
        gofmt -s $f | diff $f -
    fi
done

exit $E_VAL