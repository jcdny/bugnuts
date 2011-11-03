#!/bin/bash
ROOT=~/src/ai/bot

BOT=bugnutsv6
REMOTENAME=bugnutsv5

BHOST=ants.fluxid.pl

ARCH="`uname -s`"
DATE="`date +%Y%m%d-%H%M`"

LOG=$ROOT/log/tcp/log.${BHOST}.$DATE
ERR=$ROOT/log/tcp/err.${BHOST}.$DATE

EXE=$ROOT/bin/$ARCH/${BOT}

if [ ! -x $EXE ]; then 
    echo "No executable found at $EXE"
    exit 1
fi

while [ 1 ] ; do 
    echo "INFO: $EXE started at `date` LOG $LOG"
    echo "python $ROOT/bin/tcpclient.py $BHOST 2081 $EXE $BOT donuts -1 > $LOG 2> $ERR"
    python $ROOT/bin/tcpclient.py $BHOST 2081 $EXE $RBOT donuts -1 >> $LOG 2>> $ERR
    echo "INFO: failed at `date`, sleeping 15min" >> $LOG
    echo "INFO: failed at `date`"
    sleep 900
done
