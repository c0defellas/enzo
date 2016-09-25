all: install test

test:
	go test -v ./...

install:
	go install ./...
	@du -h $(GOPATH)/bin/cat
	strip -s $(GOPATH)/bin/cat
	@du -h $(GOPATH)/bin/cat
