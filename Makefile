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
	# go get -u github.com/swaggo/swag/cmd/swag
	# pin 1.6.7 for now. something is wrong in 1.6.8
	#	https://github.com/swaggo/swag/issues/810
	go get github.com/swaggo/swag/cmd/swag@df209afeed2334a97c83aff34ea7abcad85c31f6

swagger:
	swag init -g api/router.go

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

docker-image:
	docker build -t sconwar:local .
	docker build -t sconwar:frontend -f  frontend/Dockerfile.frontend frontend/

integrity:
	cd $(BIN_DIR) && shasum *

test:
	go test  ./...
