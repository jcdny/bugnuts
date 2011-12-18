#!/bin/bash

NOW=`date +%Y%m%d-%H%M`
DEST=/tmp/dist.$NOW
rm -rf /tmp/XX
mkdir -p $DEST/tmp
mkdir -p $DEST/build

cp `find . -name \*.go \
  -a ! -name \*_test\*.go \
  -a ! -path \*/replay/\* \
  -a ! -path \*/engine/\* \
  -a ! -path \*/statbot/\* \
  -a ! -path \*/bot[3456]/\* \
  -print` $DEST/tmp
cp ../cmd/bot/main.go $DEST/tmp
cd $DEST/tmp || exit 1
for a in *.go; do
    sed -e 's:[._] "bugnuts/[^"]*"::' -e 's/^[ ]*package .*/package main/' < $a > ../$a
    rm $a
done
cd $DEST
zip -q $DEST.zip *.go && echo "$DEST.zip created"
cd $DEST/build
unzip -q $DEST.zip

echo "Building in $DEST/build"
echo "Files with bugnuts/ : `grep -l "bugnuts/" *.go`"

6g -o _go_.6 *.go
6l -o bot _go_.6

if [ ! -x ./bot ]; then
    echo "Build failed"
    exit 1
fi

./bot < ~/bot/src/cmd/bot/testdata/test.input > test.out 2> test.err || echo "WARNING: execute failed for bot"

if [ ! -s test.out ]; then
    echo "WARNING: test failed, empty moves file $DEST/build/test.out"
    exit 1
else
    BAD=`egrep -v '^(o [0-9][0-9]* [0-9][0-9]* [news]|go)$' test.out | wc -l`
    if [ "$BAD" -ne "0" ]; then
        echo "WARNING: $BAD non move lines in output"
    fi
fi

if [ -s test.err ]; then
    echo "WARNING: nonzero stderr `wc -l <$DEST/build/test.err` lines in $DEST/build/test.err"
    head -n 5 $DEST/build/test.err
fi
cp $DEST.zip ~/tmp
