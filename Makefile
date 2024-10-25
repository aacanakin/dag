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

ci: build test.coverage

build:
	$(GOBUILD) ./...

clean.testcache:
	${GOCLEAN} -testcache

test.nocache: clean.testcache
	$(GOTESTALL)

test:
	$(GOTESTALL)

test.cover:
	$(GOTESTALL) -v -cover

test.cover.report: clean.testcache
	$(GOTESTALL) -v -coverprofile=coverage.out
	$(GOTOOL) cover -html=coverage.out

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
