

all: *.go
	go build -buildmode=plugin

run: *.so
	ccat -m plugin:wc.so:M
