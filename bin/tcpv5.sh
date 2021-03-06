#!/bin/bash
ROOT=~/src/ai/bot
BOT=bugnutsv6
RBOT=bugnutsv5
LHOST=127.0.0.1
BHOST=ants.fluxid.pl
ARCH="`uname -s`"
DATE="`date +%Y%m%d-%H%M`"
LOG=$ROOT/log/tcp/${BHOST}.log.$DATE
ERR=$ROOT/log/tcp/${BHOST}.err.$DATE
EXE=$ROOT/bin/$ARCH/$BOT
if [ ! -x $EXE ]; then
    echo "No executable found for $BOT in $EXE"
    exit 1
fi
echo "python $ROOT/bin/tcpclient.py $LHOST 2081 $EXE $BOT donuts -1 > $LOG 2> $ERR"
python $ROOT/bin/tcpclient.py $LHOST 2081 "$EXE" $RBOT donuts -1 > $LOG 2> $ERR
