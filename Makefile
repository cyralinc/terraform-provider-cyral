GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=gofmt
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
TEST=$$($(GOCMD) list ./...)
PROTOC=protoc
# Get latest version (tag). It is important to notice that the following
# commands restricts the build to those git-initialized folders. Thus,
# git clone the provider and run make on the same directory and avoid
# copying the source files to directories not initialized by git.
vVERSION:=$(shell git tag --sort v:refname | tail -n1)
VERSION:=$(shell git tag --sort v:refname | tail -n1 | sed 's/v//g')

HOSTNAME=cyral.com
NAMESPACE=terraform
NAME=cyral
BINARY=terraform-provider-$(NAME)_$(vVERSION)

all: install test

install:
	$(GOFMT) -w .
	mkdir -p out/
	# Build for both MacOS and Linux
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o out/darwin_amd64/$(BINARY) .
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o out/linux_amd64/$(BINARY) .
	# Store in local registry to be used by Terraform 13 and 14
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64
	cp out/darwin_amd64/$(BINARY) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64
	cp out/linux_amd64/$(BINARY) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64

clean:
	$(GOCLEAN) -i github.com/cyralinc/terraform-provider-cyral/...
	rm -f $(BINARY)
	rm -rf ./out
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}

test:
	$(GOTEST) $(TEST) -v
