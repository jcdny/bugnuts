#!/bin/bash
cd ~/ai/data || exit 1
~/bot/tools/replay/tcpants.sh &
~/bot/tools/replay/aichallenge.sh &
~/bot/tools/replay/fluxid.sh &
wait
echo "DONE `date`"