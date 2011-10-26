#!/bin/bash
GAME=1
HOST=ash.webfactional.com
BASE=http://$HOST/replay.
ROOT=~/src/ai/bot/replay
if [ ! -d $ROOT/data ]; then
    echo "Fatal $ROOT/data does not exist"
    exit 1
fi

cd $ROOT/data
mkdir -p $HOST || exit 1
cd $HOST

D="`expr $GAME / 100`"

while [ $GAME -lt 1400 ]; do
    OUT=$D/replay.$GAME
    if [ -f $OUT]; then
        echo "INFO: found game $GAME"
    else
        echo "INFO: getting ${BASE}$GAME"
        curl --create-dirs -o $OUT ${BASE}$GAME || echo ouch
        ls -l1 $OUT
        sleep 15
    fi
    GAME="`expr $GAME + 1`"
done
