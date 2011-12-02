#!/bin/bash

NOW=`date +%Y%m%d-%H%M`
DEST=/tmp/dist.$NOW
rm -rf /tmp/XX
mkdir -p $DEST/tmp

cp `find . -name \*.go \
  -a ! -name \*_test\*.go \
  -a ! -path \*/replay/\* \
  -a ! -path \*/engine/\* \
  -a ! -path \*/bot[345]/\* \
  -print` $DEST/tmp
cp ../cmd/bot/main.go $DEST/tmp
cd $DEST/tmp || exit 1
for a in *.go; do
    sed -e 's:[._] "bugnuts/[^"]*"::' -e 's/^[ ]*package .*/package main/' < $a > ../$a
    rm $a
done
cd $DEST
echo "Building in $DEST"
echo "Files with bugnuts/ : `grep -l "bugnuts/" *.go`"

6g -o _go_.6 *.go
6l -o tmp/bot _go_.6

if [ ! -x tmp/bot ]; then
    echo "Build failed"
    exit 1
fi

./tmp/bot < ~/bot/src/cmd/bot/testdata/test.input > tmp/test.out 2> tmp/test.err || echo "WARNING: execute failed for bot"

if [ ! -s tmp/test.out ]; then
    echo "WARNING: test failed, empty moves file $DEST/tmp/test.out"
    exit 1
fi

if [ -s tmp/test.err ]; then
    echo "WARNING: nonzero stderr `wc -l <$DEST/tmp/test.err` lines in $DEST/tmp/test.err"
fi

cd /tmp
zip $DEST.zip $DEST/*.go && echo "$DEST.zip created"
