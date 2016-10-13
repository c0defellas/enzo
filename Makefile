all: install test

TARGETS = echo cat fdisk

test:
	go test -v ./...

install:
	@for cmd in $(TARGETS); do \
	rm -f $(GOPATH)/bin/$$cmd; \
	done; \
	go install ./cmd/...
	@for cmd in $(TARGETS); do \
		echo -ne "\t\t"; \
		du -h $(GOPATH)/bin/$$cmd ; \
		echo -ne "stripped:\t"; \
		strip -s $(GOPATH)/bin/$$cmd ; \
		du -h $(GOPATH)/bin/$$cmd ; \
	done
