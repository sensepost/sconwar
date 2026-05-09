# ref: https://vic.demuzere.be/articles/golang-makefile-crosscompile/

V := $(shell git rev-parse --short HEAD)
PWD := $(shell pwd)
LD_FLAGS := -ldflags="-s -w"
BIN_DIR := build

default: clean swagger darwin linux windows integrity

clean:
	$(RM) $(BIN_DIR)/sconwar*
	go clean -x

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g api/router.go --parseDependency --parseInternal

install:
	go install

deps:
	go mod download -x

darwin:
	GOOS=darwin GOARCH=amd64 go build $(LD_FLAGS) -o '$(BIN_DIR)/sconwar-$(V)-darwin-amd64'
linux:
	GOOS=linux GOARCH=amd64 go build $(LD_FLAGS) -o '$(BIN_DIR)/sconwar-$(V)-linux-amd64'
windows:
	GOOS=windows GOARCH=amd64 go build $(LD_FLAGS) -o '$(BIN_DIR)/sconwar-$(V)-windows-amd64.exe'
docker:
	go build $(LD_FLAGS) -o sconwar

docker-images:
	docker build -t sconwar-api:local .
	docker build -t sconwar-ui:local -f frontend/Dockerfile.frontend frontend/

integrity:
	cd $(BIN_DIR) && shasum *

test:
	go test  ./...
