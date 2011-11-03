#!/bin/bash
ROOT=~/src/ai/bot
BOT=bugnutsv6
RBOT=bugnutsv5
BHOST=ants.fluxid.pl
ARCH="`uname -s`"
DATE="`date +%Y%m%d-%H%M`"
LOG=$ROOT/log/tcp/${BHOST}.log.$DATE
ERR=$ROOT/log/tcp/${BHOST}.err.$DATE
EXE=$ROOT/bin/$ARCH/${BOT}
if [ ! -x $EXE ]; then 
    echo "No executable found for $BOT in $EXE"
    exit 1
fi
while [ 1 ] ; do 
    echo "INFO: $EXE started at `date` LOG $LOG"
    python $ROOT/bin/tcpclient.py $BHOST 2081 $EXE $RBOT donuts -1 >> $LOG 2>> $ERR
    echo "INFO: failed at `date`"
sleep 900
done
