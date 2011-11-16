#!/bin/bash
BIN=~/src/ai/bot/bin
ROOT=~/src/ai/bot/log/tcp

BOT=bugnutsv7.1
REMOTENAME=bugnutsv5

BHOST=$1

ARCH="`uname -s`"
DATE="`date +%Y%m%d-%H%M`"

LOG=${BHOST}.$DATE.log
ERR=${BHOST}.$DATE.err

EXE=$BIN/$ARCH/${BOT}

if [ ! -x $EXE ]; then
    echo "No executable found at $EXE"
    exit 1
fi

cd $ROOT || exit 1
if [ -e STOP.$REMOTENAME  ] ; then
    echo "WARNING: STOP.$REMOTENAME exists.  Remove it to run"
    exit 1
fi

mkdir -p $ROOT/bin/
mkdir -p $ROOT/$BHOST/$DATE


echo "Starting $REMOTENAME running $BOT on $BHOST Logging $LOG"
exec > $LOG 2> $ERR
cd $ROOT/$BHOST/$DATE
echo "INFO: $EXE started at `date` LOG $ROOT"

GAME=0
while [ ! -e $ROOT/$BHOST/STOP ] ; do
    D="`expr $GAME % 1000`"
    if [ $D -eq 0 ]; then
        if [ $GAME -gt 0 ]; then
            DATE="`date +%Y%m%d-%H%M`-$GAME"
            gzip * &
            cd ..
            mkdir $DATE
            cd $DATE || exit 1
        fi
        rm -f ../latest
        ln -sf $DATE ../latest
    fi
    echo "python $BIN/tcpclient.py $BHOST 2081 $EXE $REMOTENAME donuts 1 > $GAME.log"

    MD5="`md5 < $EXE`"
    if [ ! -e  $ROOT/bin/$BOT.$MD5 ]; then
        cp $EXE $ROOT/bin/${BOT}-$MD5
    fi
    echo "INFO: executable $EXE md5 $MD5" > $GAME.log

    python $BIN/tcpclient.py $BHOST 2081 $EXE $REMOTENAME donuts 1 >> $GAME.log 2>> $GAME.err || sleep 10
    ln -sf $DATE/$GAME.log ../${REMOTENAME}.log
    ln -sf $DATE/$GAME.err ../${REMOTENAME}.err
    GAME="`expr $GAME + 1`"
done
echo "INFO: exited at `date`"
