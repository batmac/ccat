#! /bin/sh

set -x

#! /bin/sh

set -x
VERSION="dev"
COMMIT="none"
DATE=$(date +%Y-%m-%d@%H:%M:%S)
BUILTBY="build.minimal.sh"
VARS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE -X main.builtBy=$BUILTBY"


CGO_ENABLED=0 go build -v -ldflags "-s -w $VARS" -tags fileonly,nohl,nomd .

