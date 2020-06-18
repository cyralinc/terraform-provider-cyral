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
	go get -d golang.org/x/crypto@latest
	go run $(GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,localhost --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
