#!/usr/bin/env sh
ROOT=../tools
ARCH="`uname -s`"

export MAP1="testdata/maps/random_walk_04p_02.map" # translation  symmetry
export MAP2="testdata/maps/maze_04p_02.map" # annoying map
export MAP3="testdata/maps/maze_04p_01.map" # straight across from another hill.
export MAP4="testdata/maps/mmaze_04p_01.map" # 6 hills per player - insane.
export MAP5="testdata/maps/random_walk_04p_01.map" # empty
export MAP6="testdata/maps/mmaze_04p_02.map"  # 2 hills each

export MAP=$MAP2

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
    --food_rate 5 11 \
    --food_turn 19 37 \
    --food_start 75 175 \
    --food_visible 3 5 \
    --cutoff_turn 150 \
    --cutoff_percent 0.85 \
    "$@" \
    "./bot.sh" \
    "./bin/$ARCH/bugnutsv5" \
    "python $ROOT/sample_bots/python/GreedyBot.py" "python $ROOT/sample_bots/python/HunterBot.py" \




