# Define variables
GO=go
GOMOD=$(GO) mod
GOTOOL=$(GO) tool
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOFMT=$(GO) fmt 
GOCLEAN=$(GO) clean 
BINARY_NAME=delio-mini
CPU_PROF_NAME=cpu.prof
MEM_PROF_NAME=mem.prof

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
	./${BINARY_NAME}
	make clean

.PHONY: bench
bench:
	mkdir ./bench
	$(GOTEST) -cpuprofile ./bench/${CPU_PROF_NAME} -memprofile ./bench/${MEM_PROF_NAME} -benchmem -bench . | tee | grep -v Ch
	$(GOTOOL) pprof -top ./bench/${CPU_PROF_NAME}
	$(GOTOOL) pprof -top ./bench/${MEM_PROF_NAME}
	rm -rf bench