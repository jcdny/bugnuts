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


while [ $GAME -lt 7500 ]; do
    D="`expr $GAME / 1000`"
        mkdir -p $D
    DEST="$D/replay.$GAME"
    SOURCE=${BASE}.$GAME
    if [ -f $D/replay.$GAME ]; then
        echo "INFO: found game $DEST"
    else
        echo "INFO: getting $SOURCE"
        echo curl --create-dirs -o $D/replay.$GAME ${BASE}.$GAME || echo ouch
        ls -1l $DEST
        sleep 5
    fi
    GAME="`expr $GAME + 1`"
done
