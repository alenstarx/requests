.PHONY: all deps test clean
all: test
	go build -v 

deps:
	go get -v ./...

test: deps
	go test -v ./...

clean:
	go clean
