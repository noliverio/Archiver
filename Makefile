GOCMD=go
GOBUILD=$(GOCMD) build 
BUILDPATH=src/main.go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=archiver
WINDOWS_BINARY_NAME=archiver.exe
build:
	$(GOBUILD) -o $(BINARY_NAME) $(BUILDPATH)
install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(BUILDPATH)
debug_build:
	$(GOBUILD) -o $(BINARY_NAME) -gcflags="all=-N -l" $(BUILDPATH)
