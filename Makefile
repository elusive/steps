BINARY_NAME=steps

deploy/steps: build
	@xcopy ${BINARY_NAME}.exe e:\\AnyFolderName\\${BINARY_NAME}.exe /Y
	cd /e/AnyFolderName/; ./${BINARY_NAME}.exe -verbose

test:
	@go test ./... -v

build:
	@GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}.exe cmd/steps/main.go
.PHONY: build

build/linux:
	@GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux cmd/steps/main.go
	@chmod +x ${BINARY_NAME}-linux
.PHONY: build/linux

run: 
	@go run ./cmd/steps/main.go
.PHONY: run

clean: 
	@go clean
	@rm -f ${BINARY_NAME}.exe
	@rm -f steps/tmp*
	@rm -f __debug_bin.exe

# Need to mark the test action phony as we have 
# a directory with that exact name that we are 
# not intending to refer to with our "test" def.
.PHONY: test
