GO=go
GOPATH=$(shell pwd)
GOFLAGS=build
BIN=fetch

all : $(BIN)

fetch : fetch.go
	GOPATH=$(GOPATH) $(GO) $(GOFLAGS) -o $@ .

clean :
	-@ rm -rf $(BIN)
