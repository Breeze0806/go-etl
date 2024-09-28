IBM_DB_HOME := ${GOPATH}/src/github.com/ibmdb/clidriver
ifdef DB2_HOME
	IBM_DB_HOME=${DB2_HOME}
endif
export GO15VENDOREXPERIMENT=1
export GO111MODULE=on
export CGO_CFLAGS=-I${IBM_DB_HOME}/include
export CGO_LDFLAGS=-L${IBM_DB_HOME}/lib
export LD_LIBRARY_PATH=${IBM_DB_HOME}/lib	

.PHONY: all
all: lint release test

.PHONY: dependencies
dependencies:
	@echo "Installing db2 lib..."
	git clone -b v0.4.5 --depth=1 https://github.com/ibmdb/go_ibm_db $(go env GOPATH)/src/github.com/ibmdb/go_ibm_db
	cd $(go env GOPATH)/src/github.com/ibmdb/go_ibm_db/installer && go run setup.go
	@echo "Installing test dependencies..."
	go mod tidy

.PHONY: lint
lint:
	@rm -rf lint.log
	@echo "Checking format..."
	@find . -name '*.go' | xargs gofmt -w -s
	@echo "Checking vet..."
	@go vet ./... 2>&1 | tee -a lint.log
	@echo "Checking license..."
	@go run tools/license/main.go -c

.PHONY: test
test:
	@go test ./...

.PHONY: cover
cover:
	@sh cover.sh

.PHONY: release
release:
	@go generate ./...
	@cd cmd/datax && go build -ldflags="-s -w" && cd ../..
	@go run tools/datax/release/main.go
.PHONY: doc
doc:
	@godoc -http=:6080
