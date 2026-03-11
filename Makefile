.PHONY: build install test clean

BINARY := llm-gate
BUILD_DIR := ./build

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/llm-gate/

install: build
	cp $(BUILD_DIR)/$(BINARY) $(GOPATH)/bin/$(BINARY) 2>/dev/null || \
	cp $(BUILD_DIR)/$(BINARY) $(HOME)/go/bin/$(BINARY) 2>/dev/null || \
	sudo cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)
