.PHONY: all deps test clean

# uncomment under line to custom your GOPATH
# GOPATH=$(CURDIR)/.gopath

all: test
	go build -v 

deps:
	go get github.com/gin-gonic/gin
	go get -v ./...

test: deps
	go test -v ./...

clean:
	go clean
