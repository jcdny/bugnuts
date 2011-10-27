#!/bin/bash
# Doh! put dir expr in wrong place in snarf script.
# Fix all files; keep around in case I want to 
# remap dirs...
G=100
D=0
while [ $G -lt 1500 ]; do
    ND="`expr $G / 1000`"
    if [ "$ND" != "0" ]; then
        FROM=0/replay.$G
        TO=$ND/replay.$G

        if [ "$D" != "$ND" ]; then
            mkdir -p $ND
            D=$ND
        fi
        if [ -f $FROM ]; then
            mv $FROM $TO
        fi
    fi

    G="`expr $G + 1`"
done
