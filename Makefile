GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=gofmt
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
PROTOC=protoc
VERSION=$(shell cat VERSION)

all: install test

install:
	$(GOFMT) -w .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o out/darwin_amd64/terraform-provider-cyral_v$(VERSION) .
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o out/linux_amd64/terraform-provider-cyral_v$(VERSION) .

clean:
	$(GOCLEAN) -i github.com/cyralinc/terraform-provider-cyral/...
	rm -rf ./out

test:
