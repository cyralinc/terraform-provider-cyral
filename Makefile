GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=gofmt
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
PROTOC=protoc

all: install test

install:
	$(GOFMT) -w .
	$(GOBUILD) -o terraform-provider-cyral .

clean:
	$(GOCLEAN) -i github.com/cyralinc/terraform-provider-cyral/...

test: