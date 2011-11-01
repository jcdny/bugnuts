#!/bin/bash
GAME=25000
END=39000
HOST=aichallenge.org
BASE=http://aichallenge.org/game/0/
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

    DEST=$D/$GAME.replaygz
    SOURCE=${BASE}$D/$GAME.replaygz
    if [ -f $DEST ]; then
        echo "INFO: found game $DEST"
    else
        echo "INFO: getting $SOURCE"
        curl --create-dirs -o $DEST $SOURCE || echo Ouch
        ls -1l $DEST | grep gz
        sleep `expr $RANDOM % 10 + 1`
    fi
    GAME="`expr $GAME + 1`"
done
