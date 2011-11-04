#!/bin/bash
HOST=aichallenge.org
URLBASE=http://aichallenge.org/game/0/
ROOT=~/src/ai/bot/replay
SUFFIX=".replaygz"
PREFIX=""
LAST=$ROOT/data/LAST.$HOST

if [ ! -d $ROOT/data ]; then
    echo "Fatal $ROOT/data does not exist"
    exit 1
fi

cd $ROOT/data

mkdir -p $HOST || exit 1

GAME="`cat $LAST || echo 1`"
END=$1

if [ "$END" = "" -o "$GAME" = "" ]; then
    echo "FATAL: game range unkown end: $END start: $GAME"
    exit 1
fi

DATE="`date +%Y%m%d-%H%M`"
LOG="log/log.${HOST}.$DATE"
ERR="log/err.${HOST}.$DATE"
echo "INFO: getting $HOST games $GAME to $END log $LOG"
exec > $LOG 2> $ERR

cd $HOST

while [ $GAME -lt $END ]; do
    D="`expr $GAME / 1000`"
    mkdir -p $D
    FILENAME="$PREFIX$GAME$SUFFIX"
    DEST="$D/$FILENAME"
    SOURCE=${URLBASE}$DEST
    if [ -f $DEST ]; then
        echo "INFO: found game $DEST"
    else
        echo "INFO: getting $SOURCE to `pwd`/$DEST"
        curl --create-dirs -o $DEST $SOURCE || echo Ouch
        ls -1l $DEST
        sleep `expr $RANDOM % 10 + 10`
    fi
    GAME="`expr $GAME + 1`"
    echo $GAME > $LAST
done
