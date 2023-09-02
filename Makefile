BINARY_NAME=steps
PKG = github.com/elusive/steps
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`


# build the ldFlags parameters with version and build values
FLAGS = -ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

deploy/steps: build
	@xcopy ${BINARY_NAME}.exe e:\\AnyFolderName\\${BINARY_NAME}.exe /Y
	cd /e/AnyFolderName/; ./${BINARY_NAME}.exe -verbose

test:
	@go test ./... -v

build:
	@GOARCH=amd64 GOOS=windows go build ${FLAGS} -o ${BINARY_NAME}.exe cmd/steps/main.go

build/linux:
	@GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux cmd/steps/main.go
	@chmod +x ${BINARY_NAME}-linux

run: build
	@./steps.exe -verbose
	

clean: 
	@go clean
	@rm -f ${BINARY_NAME}.exe
	@rm -f steps-linux
	@rm -f steps/tmp*
	@rm -f __debug_bin.exe

# Need to mark the test action phony as we have 
# a directory with that exact name that we are 
# not intending to refer to with our "test" def.
.PHONY: test build
