#!/bin/bash

find ./ -type f \
    | while read i; do
    git blame -t $i 2>/dev/null;
    done \
    | sed 's/^[0-9a-f]\{8\} [^(]*(\([^)]*\) [-+0-9 ]\{14,\}).*/\1/;s/ *$//' \
    | awk '{a[$0]++; t++} END{for(n in a) if (a[n]*100.0/t > 0.5) print n}' \
    | grep -v '^breeze$' \
    | grep -v '^Breeze0808$' \
    | grep -v '@' \
    | sort -u