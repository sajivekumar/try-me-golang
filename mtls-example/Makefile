# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=server
BINARY_UNIX=$(BINARY_NAME)

# All target
all: test buildbin

# Test target
test:
	$(GOTEST) -v ./...

# Build target
buildbin:
	$(GOCMD) mod tidy
	$(GOCMD) get -v ./...
	$(GOBUILD) -o ./build/$(BINARY_NAME) ./

# Clean target
clean:
	$(GOCLEAN)
	rm -f ./build/$(BINARY_NAME)
	rm -f ./build/$(BINARY_UNIX)

# Run target
run:
	$(GOBUILD) -o ./build/$(BINARY_NAME) ./cmd/api-server
	./build/$(BINARY_NAME)

