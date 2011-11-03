#!/usr/bin/env sh
ROOT=../tools
ARCH="`uname -s`"

#export MAP="testdata/maps/maze_02p_02.map"
export MAP="testdata/maps/maze_04p_01.map"
#export MAP="testdata/maps/big"
#    --nolaunch \

../tools/playgame.py \
    --fill \
    -I -E -e -O -R \
    --player_seed 43 \
    --engine_seed 43 \
    --end_wait=0.25 \
    --verbose \
    --log_dir log \
    --turns 1000 \
    --map_file $MAP \
    "$@" \
    "./bot.sh" \
    "./bin/$ARCH/bugnutsv5" \
    "python $ROOT/sample_bots/python/GreedyBot.py" "python $ROOT/sample_bots/python/HunterBot.py" \

