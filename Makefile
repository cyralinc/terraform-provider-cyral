GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=gofmt
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
PROTOC=protoc
# Get latest version (tag)
VERSION:=$(shell git tag --sort v:refname | tail -n1)

all: install test

install:
	$(GOFMT) -w .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o out/darwin_amd64/terraform-provider-cyral_$(VERSION) .
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o out/linux_amd64/terraform-provider-cyral_$(VERSION) .

docker-go-build:
	docker-compose run app $(GOFMT) -w .
	docker-compose run -e GOOS=darwin -e GOARCH=amd64 app $(GOBUILD) -o out/darwin_amd64/terraform-provider-cyral_$(VERSION) .
	docker-compose run -e GOOS=linux -e GOARCH=amd64 app $(GOBUILD) -o out/linux_amd64/terraform-provider-cyral_$(VERSION) .

clean:
	$(GOCLEAN) -i github.com/cyralinc/terraform-provider-cyral/...
	rm -rf ./out

docker-go-clean:
	docker-compose run app $(GOCLEAN) -i github.com/cyralinc/terraform-provider-cyral/...
	rm -rf ./out

test:
	$(GOTEST) ./cyral
