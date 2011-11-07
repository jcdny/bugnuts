#!/bin/bash
HOST=ants.fluxid.pl
URLBASE=http://ants.fluxid.pl/replay
ROOT=~/src/ai/bot/replay
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
echo "INFO: getting $HOST games $GAME to $END"

cd $HOST

while [ $GAME -lt $END ]; do
    D="`expr $GAME / 1000`"
        mkdir -p $D
    DEST="$D/replay.$GAME"
    SOURCE=${URLBASE}.$GAME
    if [ -f $D/replay.$GAME ]; then
        echo "INFO: found game $DEST"
    else
        echo "INFO: getting $SOURCE"
        curl --create-dirs -o $D/replay.$GAME ${URLBASE}.$GAME || echo Ouch
        if [ -s $DEST ]; then
            chmod 444 $DEST
            ls -1l $DEST
        fi
        sleep `expr $RANDOM % 10 + 10`
    fi
    GAME="`expr $GAME + 1`"
    echo $GAME > $LAST
done
