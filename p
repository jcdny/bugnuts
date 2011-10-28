#!/usr/bin/env sh
ROOT=../tools
ARCH="`arch`"
../tools/playgame.py \
    -I -E -O -R \
    --player_seed 42 \
    --engine_seed 42 \
    --end_wait=0.25 \
    --verbose \
    --log_dir log \
    --turns 1000 \
    --map_file testdata/maps/maze_04p_01.map \
    "./MyBot" \
    "./bin/$ARCH/bugnutsv3" \
    "python $ROOT/sample_bots/python/HunterBot.py" \
    "python $ROOT/sample_bots/python/GreedyBot.py"
