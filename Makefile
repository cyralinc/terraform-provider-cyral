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
VERSION+sha:=$(VERSION)+$(shell git rev-parse --short HEAD)

HOSTNAME=cyral.com
NAMESPACE=terraform
NAME=cyral
BINARY=terraform-provider-$(NAME)_$(vVERSION)

all: local/clean local/install local/test

local/build:
	$(GOFMT) -w .
	mkdir -p out/
	# Build for both MacOS and Linux
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o out/darwin_amd64/$(BINARY) .
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o out/linux_amd64/$(BINARY) .

local/install: local/build
# Store in local registry to be used by Terraform 13 and 14
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64
	cp out/darwin_amd64/$(BINARY) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64
	cp out/linux_amd64/$(BINARY) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64

docker/build:
	docker-compose run app $(GOFMT) -w .
	docker-compose run -e GOOS=darwin -e GOARCH=amd64 app $(GOBUILD) -o out/darwin_amd64/terraform-provider-cyral_v$(VERSION) .
	docker-compose run -e GOOS=linux -e GOARCH=amd64 app $(GOBUILD) -o out/linux_amd64/terraform-provider-cyral_v$(VERSION) .
	# Store in local registry to be used by Terraform 13 and 14
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64
	cp out/darwin_amd64/$(BINARY) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/darwin_amd64
	cp out/linux_amd64/$(BINARY) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/linux_amd64

clean: local/clean
local/clean:
	$(GOCLEAN) -i github.com/cyralinc/terraform-provider-cyral/...
	rm -f $(BINARY)
	rm -rf ./out
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}

docker/clean:
	docker-compose run app $(GOCLEAN) -i github.com/cyralinc/terraform-provider-cyral/...
	rm -rf ./out
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}

local/test:
	$(GOTEST) $(TEST) -v -race

docker-compose/build: docker-compose/lint
	docker-compose build --build-arg VERSION="$(VERSION+sha)" build

docker-compose/lint:
	docker-compose run lint
