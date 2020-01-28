export GO111MODULE = on
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_PATH= "./bin"
BINARY_NAME= kubeall
BINARY_UNIX=$(BINARY_NAME)_unix
     

.DEFAULT_GOAL := all
all: test build
build:
	$(GOBUILD) -o $(BINARY_PATH)/$(BINARY_NAME) -v
test:  
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN) 
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/verchol/kubectx
         