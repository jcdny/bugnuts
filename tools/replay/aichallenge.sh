#!/bin/bash
HOST=aichallenge.org
URLBASE=http://aichallenge.org/game/0/
ROOT=~/ai/
SUFFIX=".replaygz"
PREFIX=""
LAST=$ROOT/data/$HOST/LAST
LOCK=$ROOT/data/$HOST/LOCK
STOP=$ROOT/data/$HOST/STOP

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
echo "$$ `date`" > $HOST/LOCK

GAME="`cat $LAST || echo 1`"
END=$(curl -s http://aichallenge.org/games.php | egrep 'visualizer.php.*game=[0-9]' | head -n1 | sed 's/.*game=\([0-9]*\)[^0-9].*/\1/')

if [ "$END" = "" -o "$GAME" = "" ]; then
    echo "FATAL: game range error end: \"$END\" start: \"$GAME\""
    exit 1
fi
if [ "$END" -lt "$GAME" ]; then
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

while [ $GAME -lt $END -a ! -f $STOP ]; do
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
        if [ -s $DEST ]; then
            chmod 444 $DEST
            ls -1l $DEST
        fi
        curl -o /dev/null -d "email=jeff.davis@gmail.com&stat=$HOST:game&value=$GAME" http://api.stathat.com/ez &
        curl -o /dev/null -d "email=jeff.davis@gmail.com&stat=$HOST:downloaded&count=1" http://api.stathat.com/ez &
        sleep `expr $RANDOM % 5 + 2`
    fi
    GAME="`expr $GAME + 1`"
    echo $GAME > $LAST
done
rm -f $LOCK

