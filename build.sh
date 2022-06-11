#! /bin/sh

set -x

CURLDIR="/opt/homebrew/opt/curl"

export CGO_LDFLAGS="-L $CURLDIR/lib/"
export CGO_CPPFLAGS="-I $CURLDIR/include/curl/"
go install -v  -ldflags '-w ' -tags libcurl,crappy
go build -v  -ldflags '-w ' -tags libcurl,crappy .

