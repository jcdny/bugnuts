#!/bin/bash
ROOT=~/src/ai/bot

BOT=bugnutsv7.1
REMOTENAME=bugnutsv5

BHOST=ants.fluxid.pl

ARCH="`uname -s`"
DATE="`date +%Y%m%d-%H%M`"

LOG=${BHOST}.$DATE.log
ERR=${BHOST}.$DATE.err

EXE=$ROOT/bin/$ARCH/${BOT}

if [ ! -x $EXE ]; then
    echo "No executable found at $EXE"
    exit 1
fi

cd $ROOT/log/tcp || exit 1
mkdir -p $ROOT/log/tcp/$DATE
ln -sf latest $DATE
cd $DATE || exit 1
echo "Starting $REMOTENAME running $BOT on $BHOST Logging $LOG"

if [ -e ../STOP.$REMOTENAME  ] ; then
    echo "WARNING: STOP.$REMOTENAME exists.  Remove it to run"
    exit 1
fi

exec > ../$LOG 2> ../$ERR
GAME=0
echo "INFO: $EXE started at `date` LOG $ROOT/log/tcp"
while [ ! -e ../STOP.$REMOTENAME ] ; do
    echo "python $ROOT/bin/tcpclient.py $BHOST 2081 $EXE $BOT donuts 1 > $GAME.log"
    python $ROOT/bin/tcpclient.py $BHOST 2081 $EXE $REMOTENAME donuts 1 > $GAME.log 2> $GAME.err || sleep 10
    ln -sf $GAME.log ../${REMOTENAME}.log
    ln -sf $GAME.err ../${REMOTENAME}.err
    GAME="`expr $GAME + 1`"
done
echo "INFO: exited at `date`"
