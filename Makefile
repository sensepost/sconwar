# ref: https://vic.demuzere.be/articles/golang-makefile-crosscompile/

V := $(shell git rev-parse --short HEAD)
PWD := $(shell pwd)
LD_FLAGS := -ldflags="-s -w"
BIN_DIR := build

default: clean swagger darwin linux windows integrity

clean:
	$(RM) $(BIN_DIR)/sconwar*
	go clean -x

swagger-install:
	go get -u github.com/swaggo/swag/cmd/swag

swagger:
	swag init -g api/router.go

install:
	go install

darwin:
	GOOS=darwin GOARCH=amd64 go build $(LD_FLAGS) -o '$(BIN_DIR)/sconwar-$(V)-darwin-amd64'
linux:
	GOOS=linux GOARCH=amd64 go build $(LD_FLAGS) -o '$(BIN_DIR)/sconwar-$(V)-linux-amd64'
windows:
	GOOS=windows GOARCH=amd64 go build $(LD_FLAGS) -o '$(BIN_DIR)/sconwar-$(V)-windows-amd64.exe'

integrity:
	cd $(BIN_DIR) && shasum *

test:
	go test  ./...
