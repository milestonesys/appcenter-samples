# Folders to store temporary files in
OBJ_DIR = obj
BIN_DIR = bin

# Default target name
TARGET := $(BIN_DIR)/apigateway-webserver

# By default we build for amd64 (x86-64)
ARCH ?= amd64

# By default we build for linux
GOOS ?= linux

# Builds the application generating an executable for a linux OS
build: clean
	@GOARCH=$(ARCH) go get -x -v ./src/cmd/ && \
	 GOARCH=$(ARCH) GOOS=$(GOOS) CGO_ENABLED=0 go build -o $(TARGET) ./src/cmd/

# Will locate and delete the generated executable
clean:
	@mkdir -p $(BIN_DIR) && \
	 go clean && \
	 rm -fr ${BIN_DIR} && \
	 rm -fr ${OBJ_DIR}

# Runs the application as a subprocess of the current terminal
run: build
	@./$(BIN_DIR)/apigateway-webserver
