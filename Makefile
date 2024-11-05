# Base golang tooling
GOCMD=go
GORUN=$(GOCMD) run main.go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOTESTALL=$(GOTEST) -v ./...
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
GOLINT=golangci-lint


ci: lint build test.cover

build:
	$(GOBUILD) ./...

clean.testcache:
	$(GOCLEAN) -testcache

install.linter:
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: install.linter
	$(GOLINT) run

test: clean.testcache
	$(GOTESTALL)

test.cover: clean.testcache
	$(GOTESTALL) -v -cover

test.cover.report: clean.testcache
	$(GOTESTALL) -v -coverprofile=coverage.out
	$(GOTOOL) cover -html=coverage.out

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
