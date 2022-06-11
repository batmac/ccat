#! /bin/sh

set -x

CGO_ENABLED=0 go build -v -ldflags '-w ' -tags fileonly,nohl .

