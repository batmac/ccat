#! /usr/bin/env bash

CCAT=./ccat
CAT=/bin/cat
CKSUM=md5sum
DIR=$1

for FILE in $DIR/* ; do
	echo testing $FILE
	EXPECTED=$($CKSUM $FILE| cut -d" " -f1)
	echo expected checksum is $EXPECTED
	for a in gzip lz4 lzma2 lzma s2 snap xz zlib zip zstd ; do
		echo testing $a
		SUM=`$CAT $FILE| $CCAT -m $a | $CCAT -m un$a | $CKSUM | cut -d" " -f1`
		echo $SUM
		[ "$SUM" = "$EXPECTED" ] || (echo FAILED && exit 1)
	done
done
