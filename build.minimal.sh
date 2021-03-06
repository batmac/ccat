#! /bin/sh

set -x

GIT=$(git describe --tags 2> /dev/null || git rev-parse --short HEAD)
VERSION="$GIT"
COMMIT=$(git rev-parse HEAD)
DATE=$(date +%Y-%m-%d@%H:%M:%S)
BUILTBY="build.minimal.sh"
TAGS=fileonly,nohl,nomd
VARS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE -X main.builtBy=$BUILTBY -X main.tags=$TAGS"

go test -tags $TAGS ./...
cd cmd/ccat || exit 1

CGO_ENABLED=0 go build -v -ldflags "-s -w $VARS" -tags $TAGS .

mv ccat ../../ccat
