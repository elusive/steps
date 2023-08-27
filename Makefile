BINARY_NAME=steps

run/steps:
	@go run ./cmd/steps/main.go

test:
	@go test ./... -v

build:
	@GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}.exe cmd/steps/main.go

run: build
	@./${BINARY_NAME} -verbose

clean: 
	@go clean
	@rm -f ${BINARY_NAME}.exe
	@rm -f steps/tmp*
	@rm -f __debug_bin.exe

# Need to mark the test action phony as we have 
# a directory with that exact name that we are 
# not intending to refer to with our "test" def.
.PHONY: test
