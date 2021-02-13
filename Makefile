GO ?= go
GOFMT ?= gofmt -s
GOLINT ?= ~/go/bin/golint
GOFILES := $(shell find . -name "*.go")
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /examples/)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

vet:
	$(GO) vet $(VETPACKAGES)

lint:
	$(GOLINT) $(GOFILES)

v1test:
	$(GO) test -v ./...

cover:
	$(GO) test github.com/yomorun/y3-codec-golang/pkg -coverprofile=prof.out && $(GO) tool cover -html=prof.out && rm prof.out

test:
	$(GO) test -v api.go api_test.go stream_api.go stream_api_test.go

test-spec:
	$(GO) test -v github.com/yomorun/y3-codec-golang/pkg/spec
