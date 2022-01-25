GOTOOLFLAGS = GO15VENDOREXPERIMENT=1
GO ?= $(GOTOOLFLAGS) go
GOLINT ?= golint

GORELOS = linux
GORELARCH = amd64


all:: fmt build 

build: 
	$(GO) build -o ./bin ./...

fmt:
	$(GO) fmt ./...

