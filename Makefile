# Define variables
GO=go
GOMOD=$(GO) mod
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOFMT=$(GO) fmt 
GOCLEAN=$(GO) clean 
BINARY_NAME=delio-mini

# Define targets
.PHONY: all
all: build

.PHONY: deps
deps:
	$(GOCLEAN) 
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: test
test: deps
	$(GOTEST) ./...

.PHONY: build
build: test
	$(GOBUILD) -o $(BINARY_NAME) .

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: fmt
fmt:
	$(GOFMT) ./...

.PHONY: gh-workflow
gh-workflow:
	set -e
	make all
	./delio-mini
	make clean
