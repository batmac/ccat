#! /bin/sh

set -x

GIT=$(git tag|tail -n1)
VERSION=">$GIT+dev"
COMMIT=$(git rev-parse HEAD)
DATE=$(date +%Y-%m-%d@%H:%M:%S)
BUILTBY="build.sh"
TAGS=libcurl,crappy
VARS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE -X main.builtBy=$BUILTBY -X main.tags=$TAGS"

#CURLDIR="/opt/homebrew/opt/curl"
#export CGO_LDFLAGS="-L $CURLDIR/lib/"
#export CGO_CPPFLAGS="-I $CURLDIR/include/curl/"
go test -tags $TAGS ./...
cd cmd/ccat || exit 1
# go install -v  -ldflags "-s -w $VARS" -tags $TAGS
go build -v  -ldflags "-s -w $VARS" -tags $TAGS .

mv ccat ../../ccat

