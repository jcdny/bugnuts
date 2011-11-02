#!/bin/bash
GAME=7500
END=15700
HOST=ants.fluxid.pl
BASE=http://ants.fluxid.pl/replay
ROOT=~/src/ai/bot/replay
if [ ! -d $ROOT/data ]; then
    echo "Fatal $ROOT/data does not exist"
    exit 1
fi

cd $ROOT/data
mkdir -p $HOST || exit 1

DATE="`date +%Y%m%d-%H%M`"
LOG="${HOST}.log.$DATE"
ERR="${HOST}.err.$DATE"
echo "INFO: getting $HOST games $GAME to $END log $LOG"
exec > $LOG 2> $ERR

cd $HOST
while [ $GAME -lt $END ]; do
    D="`expr $GAME / 1000`"
        mkdir -p $D
    DEST="$D/replay.$GAME"
    SOURCE=${BASE}.$GAME
    if [ -f $D/replay.$GAME ]; then
        echo "INFO: found game $DEST"
    else
        echo "INFO: getting $SOURCE"
        curl --create-dirs -o $D/replay.$GAME ${BASE}.$GAME || echo Ouch
        ls -1l $DEST
        sleep `expr $RANDOM % 10 + 10`
    fi
    GAME="`expr $GAME + 1`"
done
