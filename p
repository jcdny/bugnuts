#!/usr/bin/env sh
ROOT=../tools
ARCH="`uname -s`"

export MAP="testdata/maps/random_walk_04p_02.map" # translation  symmetry
export MAP="testdata/maps/random_walk_04p_01.map" # empty
export MAP="testdata/maps/mmaze_04p_01.map" # 6 hills per player - insane.
export MAP="testdata/maps/maze_04p_02.map" # annoying map
export MAP="testdata/maps/maze_04p_01.map" # straight across from another hill.
export MAP="testdata/maps/mmaze_04p_02.map"  # 2 hills each

#export MAP="testdata/maps/big"
#    --nolaunch \

../tools/playgame.py \
    --fill \
    -I -E -e -O -R \
    --player_seed 41 \
    --engine_seed 41 \
    --end_wait=0.25 \
    --verbose \
    --log_dir log \
    --turns 1000 \
    --map_file $MAP \
    "$@" \
    "./bot.sh" \
    "./bin/$ARCH/bugnutsv5" \
    "python $ROOT/sample_bots/python/GreedyBot.py" "python $ROOT/sample_bots/python/HunterBot.py" \

