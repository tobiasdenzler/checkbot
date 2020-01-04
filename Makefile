# Binary name
BINARY=checkbot

# Build values
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

# Setup the -ldflags option for go build here
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

build:
	@echo "Version: ${VERSION}"
	@echo "Build: ${BUILD}"
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build ${LDFLAGS} -a -installsuffix cgo -o ${BINARY}_${VERSION} ./cmd/server/

run:
	go run ${LDFLAGS} ./cmd/server/ -logLevel=info -enableSandbox=true

clean:
	if [ -f ${BINARY}_v* ] ; then rm -f ${BINARY}_v* ; fi

.PHONY: all