#! /bin/sh

set -x

#! /bin/sh

set -x
GIT=$(git tag|tail -n1)
VERSION=">$GIT+dev"
COMMIT=$(git rev-parse HEAD)
DATE=$(date +%Y-%m-%d@%H:%M:%S)
BUILTBY="build.minimal.sh"
TAGS=fileonly,nohl,nomd
VARS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE -X main.builtBy=$BUILTBY -X main.tags=$TAGS"

CGO_ENABLED=0 go build -v -ldflags "-s -w $VARS" -tags $TAGS .

