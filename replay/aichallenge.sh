#!/bin/bash
GAME=1
HOST=aichallenge.org
BASE=http://aichallenge.org/game/0/
ROOT=~/src/ai/bot/replay
if [ ! -d $ROOT/data ]; then
    echo "Fatal $ROOT/data does not exist"
    exit 1
fi

cd $ROOT/data
mkdir -p $HOST || exit 1
cd $HOST

D="`expr $GAME / 100`"

while [ $GAME -lt 17500 ]; do
    OUT=$D/$GAME.replaygz
    if [ -f $OUT ]; then
        echo "INFO: found game $GAME"
    else
        echo "INFO: getting ${BASE}$D/$GAME.replaygz"
        curl --create-dirs -o $OUT ${BASE}$D/$GAME.replaygz || echo ouch
        ls -1l $OUT | grep gz
        sleep 3
    fi
    GAME="`expr $GAME + 1`"
done
