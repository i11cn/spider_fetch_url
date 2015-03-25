GO=go
GOPATH=/home/snower/golang/pkgs
GOFLAGS=build
BIN=fetch

all : $(BIN)

fetch : fetch.go
	GOPATH=$(GOPATH) $(GO) $(GOFLAGS) -o $@ .

clean :
	-@ rm -rf $(BIN)
