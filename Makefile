all: ccat readme

ccat: *.go */*.go
	./build.sh

readme: ccat
	cp README.header.md README.md
	echo '```'      >> README.md
	sh -c './ccat -h 2>&1'  >> README.md
	echo '```'      >> README.md
