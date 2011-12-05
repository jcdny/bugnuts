#!/bin/bash
HOST=ants.fluxid.pl
URLBASE=http://ants.fluxid.pl/replay
ROOT=~/ai
LAST=$ROOT/data/$HOST/LAST
LOCK=$ROOT/data/$HOST/LOCK

if [ ! -d $ROOT/data ]; then
    echo "Fatal $ROOT/data does not exist"
    exit 1
fi
cd $ROOT/data
mkdir -p $HOST || exit 1

if [ -f $LOCK ]; then
    echo "LOCK for $HOST exists.  remove $LOCK to run"
    exit 1
fi
echo `date` > $HOST/LOCK

GAME="`cat $LAST || echo 1`"
END=$(curl -s http://$HOST | egrep "href='/replay.[0-9]" | head -n1 | sed -e 's#^.*/replay\.\([0-9]*\)[^0-9].*#\1#')

if [ "$END" = "" -o "$GAME" = "" ]; then
    echo "FATAL: game range error end: \"$END\" start: \"$GAME\""
    exit 1
fi
if [ "$GAME" -gt "$END" ]; then
    echo "FATAL: end before start, End: \"$END\" Start: \"$GAME\""
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
    if [ -f $D/replay.$GAME -o -f $D/replay.${GAME}.gz ]; then
        echo "INFO: found game $DEST"
    else
        echo "INFO: getting $SOURCE"
        curl --create-dirs -o $D/replay.$GAME ${URLBASE}.$GAME || echo Ouch
        if [ -s $DEST ]; then
            gzip $DEST
            chmod 444 ${DEST}.gz
            ls -1l ${DEST}.gz
        fi
        sleep `expr $RANDOM % 10 + 10`
    fi
    GAME="`expr $GAME + 1`"
    echo $GAME > $LAST
done
rm -f $LOCK