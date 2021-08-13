FROM hashicorp/terraform:1.0.4 as terraform

FROM golang:1.14.15-alpine3.13 AS build
WORKDIR /go/src/cyral
RUN apk add --no-cache build-base=0.5-r2
COPY main.go go.mod go.sum ./
COPY client/ client/
COPY cyral/ cyral/
COPY doc/ doc/
COPY scripts/ scripts/
RUN gofmt -w . \
    && go test ./... -race \
    && mkdir -p /out \
    && GOOS=darwin GOARCH=amd64 go build -o out/darwin_amd64/terraform-provider-cyral . \
    && GOOS=linux GOARCH=amd64 go build -o out/linux_amd64/terraform-provider-cyral .

FROM alpine:3.13.5 as output
ARG VERSION
RUN mkdir -p /root/.terraform.d/plugins/local/terraform/cyral/${VERSION:?You must set the VERSION build argument}
COPY --from=build /go/src/cyral/out/ /root/.terraform.d/plugins/local/terraform/cyral/${VERSION}
COPY --from=terraform /bin/terraform /bin/terraform