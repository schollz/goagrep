SOURCEDIR=.
SOURCES := main.go

BINARY=goagrep

VERSION=2.0
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.VersionNum=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	rm -rf builds
	mkdir builds
	go get github.com/schollz/goagrep/goagrep
	go get github.com/firstrow/tcp_server
	go get github.com/codegangsta/cli
	go build ${LDFLAGS} -o builds/${BINARY} ${SOURCES}

.PHONY: install
install:
	$(MAKE) clean
	$(MAKE)
	sudo mv goagrep /usr/local/bin/
	echo "Installed to /usr/local/bin/"

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf builds
	rm -rf goagrep*

.PHONY: binaries
binaries:
	rm -rf builds
	mkdir builds
	# Build Windows
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o goagrep.exe -v *.go
	zip -r goagrep_${VERSION}_windows_amd64.zip goagrep.exe LICENSE
	mv goagrep_${VERSION}_windows_amd64.zip builds/
	rm goagrep.exe
	# Build Linux
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o goagrep -v *.go
	zip -r goagrep_${VERSION}_linux_amd64.zip goagrep LICENSE
	mv goagrep_${VERSION}_linux_amd64.zip builds/
	rm goagrep
	# Build OS X
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o goagrep -v *.go
	zip -r goagrep_${VERSION}_osx.zip goagrep LICENSE
	mv goagrep_${VERSION}_osx.zip builds/
	rm goagrep
	# Build Raspberry Pi / Chromebook
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o goagrep -v *.go
	zip -r goagrep_${VERSION}_linux_arm.zip goagrep LICENSE
	mv goagrep_${VERSION}_linux_arm.zip builds/
	rm goagrep
