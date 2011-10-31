#!/usr/bin/env sh
ROOT=../tools
ARCH="`uname -s`"

export MAP="testdata/maps/maze_04p_01.map"

#    --nolaunch \

../tools/playgame.py \
    -I -E -e -O -R \
    --player_seed 42 \
    --engine_seed 42 \
    --end_wait=0.25 \
    --verbose \
    --log_dir log \
    --turns 1000 \
    --map_file $MAP \
    "$@" \
    "./bot.sh" \
    "./bin/$ARCH/bugnutsv4" \
    "python $ROOT/sample_bots/python/HunterBot.py" \
    "python $ROOT/sample_bots/python/GreedyBot.py"
