ifndef GOPATH
$(error $$GOPATH is not set)
endif

CMD := "./cmd"
TARGETS := $(shell ls -l $(CMD) | awk '/^d/ { print $$NF }')

all: install test

test:
	go test -v ./...

install:
	@for cmd in $(TARGETS); do \
		rm -f $(GOPATH)/bin/$$cmd; \
	done; \
	go install $(CMD)/...
	@for cmd in $(TARGETS); do \
		echo -ne "\t\t"; \
		du -h $(GOPATH)/bin/$$cmd ; \
		echo -ne "stripped:\t"; \
		strip -s $(GOPATH)/bin/$$cmd ; \
		du -h $(GOPATH)/bin/$$cmd ; \
	done
