FROM hashicorp/terraform:1.10.2 as terraform

FROM golang:1.22.6-alpine3.20 AS build
WORKDIR /go/src/cyral
COPY main.go go.mod go.sum ./
COPY client/ client/
COPY cyral/ cyral/
COPY docs/ docs/
RUN gofmt -w . \
    && go test ./... -race \
    && mkdir -p /out \
    && GOOS=darwin GOARCH=amd64 go build -o out/darwin_amd64/terraform-provider-cyral . \
    && GOOS=linux GOARCH=amd64 go build -o out/linux_amd64/terraform-provider-cyral .

FROM alpine:3.20.3 as output
ARG VERSION
RUN mkdir -p /root/.terraform.d/plugins/local/terraform/cyral/${VERSION:?You must set the VERSION build argument}
COPY --from=build /go/src/cyral/out/ /root/.terraform.d/plugins/local/terraform/cyral/${VERSION}
COPY --from=terraform /bin/terraform /bin/terraform
