DB2HOME=${GOPATH}/src/github.com/ibmdb/clidriver
export GO15VENDOREXPERIMENT=1
export GO111MODULE=on
export CGO_CFLAGS=-I${DB2HOME}/include
export CGO_LDFLAGS=-L${DB2HOME}/lib
export LD_LIBRARY_PATH=${DB2HOME}/lib
# Many Go tools take file globs or directories as arguments instead of packages.
# The linting tools evolve with each Go version, so run them only on the latest
# stable release.
GO_VERSION := $(shell go version | cut -d " " -f 3)
GO_MINOR_VERSION := $(word 2,$(subst ., ,$(GO_VERSION)))
LINTABLE_MINOR_VERSIONS := 18
ifneq ($(filter $(LINTABLE_MINOR_VERSIONS),$(GO_MINOR_VERSION)),)
SHOULD_LINT := true
endif

.PHONY: all
all: lint release test

.PHONY: dependencies
dependencies:
	@echo "Installing db2 lib..."
	git clone -b v0.4.1 --depth=1 https://github.com/ibmdb/go_ibm_db ${GOPATH}/src/github.com/ibmdb/go_ibm_db
	cd ${GOPATH}/src/github.com/ibmdb/go_ibm_db/installer && go run setup.go
ifdef SHOULD_LINT
	@echo "Installing golint..."
	go install -v golang.org/x/lint/golint@latest
else
	@echo "Not installing golint, since we don't expect to lint on" $(GO_VERSION)
endif
	@echo "Installing test dependencies..."
	go mod tidy

.PHONY: lint
lint:
ifdef SHOULD_LINT
	@rm -rf lint.log
	@echo "Checking format..."
	@find . -name '*.go' | xargs gofmt -w -s
	@echo "Checking vet..."
	@go vet ./... 2>&1 | tee -a lint.log
	@echo "Checking lint..."
	@golint ./... 2>&1 | tee -a lint.log
	@[ ! -s lint.log ]
	@echo "Checking license..."
	@go run tools/license/main.go -c
else
	@echo "Skipping linters on" $(GO_VERSION)
endif

.PHONY: test
test:
	@go test ./...

.PHONY: cover
cover:
	@sh cover.sh

.PHONY: release
release:
	@go generate ./... && cd cmd/datax && go build && cd ../..

.PHONY: doc
doc:
	@godoc -http=:6080