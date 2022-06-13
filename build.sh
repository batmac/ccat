#! /bin/sh

set -x
VERSION="dev"
COMMIT="none"
DATE=$(date +%Y-%m-%d@%H:%M:%S)
BUILTBY="build.sh"
VARS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE -X main.builtBy=$BUILTBY"

CURLDIR="/opt/homebrew/opt/curl"
export CGO_LDFLAGS="-L $CURLDIR/lib/"
export CGO_CPPFLAGS="-I $CURLDIR/include/curl/"
go install -v  -ldflags "-s -w $VARS" -tags libcurl,crappy
go build -v  -ldflags "-s -w $VARS" -tags libcurl,crappy .

