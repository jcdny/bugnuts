#!/bin/bash
for a in *.log; do
    GAME="`head -n3 $a | egrep 'game [0-9]+' | sed -e 's/^[^0-9]* //' -e 's/ .*$//'`"
    if [ "$GAME" = "" ]; then 
        echo "WARNING: no game ID found $a"
    else
        
done