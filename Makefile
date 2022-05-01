export GO15VENDOREXPERIMENT=1
export GO111MODULE=on
# Many Go tools take file globs or directories as arguments instead of packages.
COVERALLS_TOKEN=477d8d1f-b729-472f-b842-e0e4b03bc0c2
# The linting tools evolve with each Go version, so run them only on the latest
# stable release.
GO_VERSION := $(shell go version | cut -d " " -f 3)
GO_MINOR_VERSION := $(word 2,$(subst ., ,$(GO_VERSION)))
LINTABLE_MINOR_VERSIONS := 16
ifneq ($(filter $(LINTABLE_MINOR_VERSIONS),$(GO_MINOR_VERSION)),)
SHOULD_LINT := true
endif

.PHONY: all
all: lint examples test

.PHONY: dependencies
dependencies:
	@echo "Installing db2 lib..."
	git clone --depth=50 https://github.com/ibmdb/go_ibm_db ${GOPATH}/src/github.com/ibmdb/go_ibm_db
	cd ${GOPATH}/src/github.com/ibmdb/go_ibm_db/installer && go run setup.go
ifdef SHOULD_LINT
	@echo "Installing golint..."
	go get -d golang.org/x/lint/golint
else
	@echo "Not installing golint, since we don't expect to lint on" $(GO_VERSION)
endif
	@echo "Installing test dependencies..."
	go mod tidy

.PHONY: lint
lint:
ifdef SHOULD_LINT
	export DB2HOME=${GOPATH}/src/github.com/ibmdb/clidriver
	export CGO_CFLAGS=-I$DB2HOME/include
	export CGO_LDFLAGS=-L$DB2HOME/lib
	export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$DB2HOME/lib
	@rm -rf lint.log
	@echo "Installing test dependencies for vet..."
	@go test ./...
	@echo "Checking vet..."
	@go vet ./... 2>&1 | tee -a lint.log
	@echo "Checking lint..."
	@golint ./... 2>&1 | tee -a lint.log
	@[ ! -s lint.log ]
else
	@echo "Skipping linters on" $(GO_VERSION)
endif

.PHONY: test
test:
	export DB2HOME=${GOPATH}/src/github.com/ibmdb/clidriver
	export CGO_CFLAGS=-I$DB2HOME/include
	export CGO_LDFLAGS=-L$DB2HOME/lib
	export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$DB2HOME/lib
	@go test ./...

.PHONY: cover
cover:
	export DB2HOME=${GOPATH}/src/github.com/ibmdb/clidriver
	export CGO_CFLAGS=-I$DB2HOME/include
	export CGO_LDFLAGS=-L$DB2HOME/lib
	export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$DB2HOME/lib
	sh cover.sh

.PHONY: examples
examples:
	export DB2HOME=${GOPATH}/src/github.com/ibmdb/clidriver
	export CGO_CFLAGS=-I$DB2HOME/include
	export CGO_LDFLAGS=-L$DB2HOME/lib
	export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$DB2HOME/lib
	@go generate ./... && cd cmd/datax && go build && cd ../..

.PHONY: doc
doc:
	@godoc -http=:6080