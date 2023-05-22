GOLANGCI_BIN=$(GOBIN)/golangci-lint
GOLANGCI_LINT_VERSION='1.51.2'

# Formatter to use; one of `gofmt`, `goimports`, or `golines`
FORMATTER := gofmt

export GO111MODULE=on
export CGO_ENABLE=1

Q=@
COVERAGE_THRESHOLD=100

# race detractor requires enabling CGO.
GOTESTFLAGS = -race
ifndef Q
GOTESTFLAGS += -v
endif

.PHONY: lint
lint:
	@bash -c '(\
		(which golangci-lint && golangci-lint --version) | grep -q $(GOLANGCI_LINT_VERSION) ||\
		curl -sfL https://raw.githubsecurecontent.com/golangci/golangci-lint/master/install.sh |\
		sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION) \
	) >\
	/dev/null && PATH="$$PATH:$(shell go env GOPATH)/bin" \
	golangci-lint run --timeout=3m --config=.golangci.yaml ./...'
	@echo Makefile target: golangci-lint run

.PHONY: vet
vet:
	$Qgo vet ./...
	@echo Makefile target: go vet

.PHONY: test
test: vet
	$QCGO_ENABLED=1 go test $(GOTESTFLAGS) -race -coverpkg=$(`go list ./... | grep -v */mocks | grep -v cmd`) -coverprofile=.coverageprofile ./...
	go tool cover -func=.coverageprofile | grep total | awk '{print substr($$3, 1, length($$3) - 1)}' | awk '{if (!($$1 >= $(COVERAGE_THRESHOLD))) { print "Coverage: " $$1 "%" ", Excepted threshold: " $(COVERAGE_THRESHOLD) "%"; exit 1} else { print "Total coverage: " $$1 }}'
	@echo Makefile target: go test

.PHONY: fmtfix
fmtfix: $(FORMATTER)
	$Q$(FORMATTER) -w $(shell find . -iname '*.go' | grep -v vendor)
	@echo Makefile target: gofmt

.PHONY: goimports
goimports:
ifeq (, $(shell which goimports))
	$QGO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	@echo Makefile target: go get goimports
endif
