SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
BIN_FOLDER=bin
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 2.0-alpha
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC := main.go 

# erase swp file
ERASE_SWP := $(shell find . -type f -name '*.swp' | xargs rm)

.PHONY: all build clean install uninstall fmt simplify check run

all: check install

build:$(ERASE_SWP)
	@go build -o $(BIN_FOLDER)/$(TARGET) ${LDFLAGS} $(SRC) 

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS)

test: 
	@go test -run ''

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: build 
	./$(BIN_FOLDER)/$(TARGET) --config ./conf.yaml 
run2: build 
	./$(BIN_FOLDER)/$(TARGET) --config ./conf.yaml --logtostderr -v 3 
run3: build 
	./$(BIN_FOLDER)/$(TARGET) --config ./conf2.yaml 

