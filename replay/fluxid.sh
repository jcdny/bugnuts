#!/bin/bash
GAME=1
HOST=ants.fluxid.pl
BASE=http://ants.fluxid.pl/replay
ROOT=~/src/ai/bot/replay
if [ ! -d $ROOT/data ]; then
    echo "Fatal $ROOT/data does not exist"
    exit 1
fi

cd $ROOT/data
mkdir -p $HOST || exit 1
cd $HOST

D="`expr $GAME / 100`"

while [ $GAME -lt 7500 ]; do
    if [ -f $D/replay.$GAME ]; then
        echo "INFO: found game $GAME"
    else
        echo "INFO: getting $GAME"
        curl --create-dirs -o $D/replay.$GAME ${BASE}.$GAME || echo ouch
        sleep 5 
    fi
    GAME="`expr $GAME + 1`"
done